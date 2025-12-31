package ApiServer

import (
	"AITranslatio/app/Model/goods"
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"strings"
	"time"
)

func (s *ApiServer) SeckillOrderHandler(ctx context.Context, body []byte) error {

	var order goods.SeckillOrder

	if err := json.Unmarshal(body, &order); err != nil {
		s.logger.Error("❌ [Consumer] 消息JSON解析失败，丢弃消息",
			zap.String("body", string(body)),
			zap.Error(err))
		return nil // 返回 nil 会导致 Msg.Ack()，消息被移除
	}

	//生产者放入MQ时会把seckill：order ：ｘｘｘ写入redis
	idempotentKey := fmt.Sprintf("seckill:done:%d", order.OrderID)
	success, err := s.redis.SetNX(ctx, idempotentKey, "1", 24*time.Hour).Result()
	if err != nil {
		return fmt.Errorf("redis err:%v", err)
	}

	if !success {
		// SetNX 返回 false，说明 Key 已存在 -> 重复消息
		s.logger.Warn("♻️ 拦截到重复消费消息", zap.Int64("order_id", order.OrderID))
		return nil // 视为成功，Ack 掉
	}

	//写数据库
	err = s.DAO.CreateSeckillOrder(order)
	if err != nil {
		s.logger.Error("落库失败", zap.Error(err))

		// 如果写库失败（非重复主键错误），必须把刚才 SetNX 的 Key 删掉
		// 否则这条消息下次重试时，会被上面的 !success 拦截，导致永远无法入库。
		if !strings.Contains(err.Error(), "Duplicate entry") {
			s.redis.Del(ctx, idempotentKey)
			return err // Nack 重试
		}
		// 如果是 Duplicate entry，说明 DB 里已经有了，也算幂等成功，不需要删 Key
		return nil
	}

	return nil
}
