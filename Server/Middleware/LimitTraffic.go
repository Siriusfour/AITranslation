package Middleware

import (
	"AITranslatio/Config/interf"
	"github.com/redis/go-redis/v9"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// LocalLimitMiddleware 创建一个基于令牌桶的本地限流中间件
// limit: 每秒产生多少令牌 (QPS)
// burst: 桶的大小 (允许瞬间并发多少)
func LocalLimitMiddleware(l *rate.Limiter) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 尝试拿一个令牌，非阻塞
		if !l.Allow() {
			// 拿不到令牌，直接返回 503 Service Unavailable
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code": 503,
				"msg":  "系统繁忙，请稍后再试 (Local Limit)",
			})
			return
		}

		// 拿到令牌，放行
		c.Next()
	}
}

func LimitByID(cfg interf.ConfigInterface, redis *redis.Client, script *redis.Script) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ip := ctx.ClientIP()

		RateLimit(ctx, redis, script, ip, cfg.GetInt("limit.limit"), cfg.GetInt("limit.windows"))
	}

}

// limit: 5 (次), window: 1 (秒)
func RateLimit(ctx *gin.Context, rdb *redis.Client, script *redis.Script, ip string, limit int, windows int) bool {
	key := "limit:ip:" + ip

	res, err := script.Run(ctx, rdb, []string{key}, limit, windows).Int()

	if err != nil {
		// 如果 Redis 挂了，通常策略是放行
		return true
	}

	return res == 1
}
