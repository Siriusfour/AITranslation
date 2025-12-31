package ApiServer

import (
	"AITranslatio/app/Model/goods"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
	"time"
)

type seckillGoods struct {
	EndTime   int64 `json:"end_time" redis:"end_time"`
	ID        int64 `json:"id" redis:"id"`
	StartTime int64 `json:"start_time" redis:"start_time"`
}

func (s *ApiServer) StartSeckill(ctx *gin.Context, userID int64) error {

	user, err := s.DAO.FindUserByID(userID, "UserID")
	if err != nil {
		return fmt.Errorf("获取用户信息失败！%w", err)
	}

	if user.Admin == false {
		return fmt.Errorf("用户权限不足")
	}

	allGoods, err := s.DAO.FindAllSeckillGoods()
	if err != nil {
		return fmt.Errorf("数据库查找失败！%w", err)
	}

	for _, good := range allGoods {
		err := PreHeatData(ctx, s.redis, good)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ApiServer) SeckillBuy(ctx *gin.Context, userID int64, goodsID, SeckillGoodsID int64) error {

	infoKey := fmt.Sprintf("seckill:info:%d", SeckillGoodsID)
	stockKey := fmt.Sprintf("seckill:stock:%d", SeckillGoodsID)
	userKey := fmt.Sprintf("seckill:users:%d", SeckillGoodsID)

	jsonStr, err := s.redis.Get(ctx, infoKey).Result()

	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("活动不存在或已结束")
	} else if err != nil {
		return fmt.Errorf("Redis读取失败: %w", err)
	}

	//根據goodsID找到该产品，判断时间
	var seckillGoods seckillGoods

	if err := json.Unmarshal([]byte(jsonStr), &seckillGoods); err != nil {
		return fmt.Errorf("数据解析失败: %w", err)
	}

	if time.Now().Unix() < seckillGoods.StartTime {
		return fmt.Errorf("活动时间未到！")
	}

	//调用redis的脚本，扣除库存
	code, err := s.scripts["Seckill"].Run(ctx, s.redis, []string{stockKey, userKey}, userID).Int()

	// 2. 先判断系统是否出错 (System Error)
	// 比如 Redis 连不上、超时、脚本写错了等
	if err != nil {
		s.logger.Error("❌ Redis 执行 Lua 失败: %w", zap.Error(err))
		return fmt.Errorf("系统繁忙，请重试")
	}
	switch code {
	case 1:
		s.logger.Info("用户 %d 抢购成功", zap.Int64("UserID", userID))
	case 0:
		return fmt.Errorf("很遗憾，商品已售罄！")
	case -1:
		return fmt.Errorf("该商品限购一件！")
	case -2:
		return fmt.Errorf("活动数据异常(未预热)！")
	default:
		return fmt.Errorf("未知状态码！ %d", code)
	}

	//如果抢货成功生成一个MQ生产者，向队列写入消息
	order := goods.SeckillOrder{
		UserID:         userID,
		OrderID:        s.snowFlakeGenerator.GetID(),
		GoodsID:        SeckillGoodsID,
		SeckillGoodsID: goodsID,
	}

	body, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("编码失败！%w", zap.Error(err))
	}

	err = s.RabbitMQ.Publish("seckill", "order", body)
	if err != nil {
		//打日志
		s.logger.Error("MQ投递失败，开始执行回滚",
			zap.Int64("uid", userID),
			zap.Int64("gid", SeckillGoodsID),
			zap.Error(err))

		s.RollbackStock(ctx, stockKey, userKey, userID) //mq发送失败，回滚事务
		return fmt.Errorf("MQ投递失败！%w", err)
	}

	return nil
}

// PreHeatData 开始秒杀活动，把mysql的秒杀产品放到redis
// redis结构：
// 商品详情 Key(Hash)                          key{seckill:info:{GoodsID}} ----value{ startTime:100000,endTime:150000 .... }
// 库存 Key(String)						    key{seckill:stock:{GoodsID}}----value 1000
// 防重（存储已经购买过的userID） Key(Set)	    key{seckill:users:{GoodsID}}----value [10086, 10087, ...]
func PreHeatData(ctx *gin.Context, redisClient *redis.Client, seckillGoods *goods.SeckillGoods) error {

	pipe := redisClient.Pipeline()
	ttl := time.Until(time.Now().Add(1 * time.Hour))

	infokey := fmt.Sprintf("seckill:info:%d", seckillGoods.SeckillID)
	userskey := fmt.Sprintf("seckill:users:%d", seckillGoods.SeckillID)
	stockkey := fmt.Sprintf("seckill:stock:%d", seckillGoods.SeckillID)

	infoData := map[string]interface{}{
		"start_time": seckillGoods.StartTime.Unix(),
		"end_time":   seckillGoods.EndTime.Unix(),
		"id":         seckillGoods.SeckillID,
	}
	infoJson, _ := json.Marshal(infoData)

	pipe.Set(ctx, stockkey, seckillGoods.StockCount, ttl)
	pipe.Set(ctx, infokey, infoJson, ttl)
	pipe.Del(ctx, userskey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis pipeline 失败: %w", err)
	}

	return nil

}

// RollbackStock 回滚逻辑
func (s *ApiServer) RollbackStock(ctx *gin.Context, stockKey, userKey string, userID int64) {
	s.redis.Incr(ctx, stockKey)
	s.redis.SRem(ctx, userKey, userID)
}
