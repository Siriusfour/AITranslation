package Middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go" // 引入 zipkin 官方库
	"go.uber.org/zap"
)

// HttpLog  来记录 HTTP Access Log
func HttpLog(logMap map[string]*zap.Logger, logKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		UserID := c.GetInt64("UserID")
		sessionID := c.GetString("SessionID")
		fmt.Println("UserID=", UserID)
		fmt.Println("SessionID=", sessionID)
		query := c.Request.URL.RawQuery
		logger, _ := logMap[logKey]
		start := time.Now()
		path := c.Request.URL.Path
		referer := c.Request.Referer()

		if path == "/metrics" {
			c.Next()
			return
		}

		c.Next()

		// --- 请求结束，开始记录日志 ---

		end := time.Now()
		latency := end.Sub(start)

		// 从 Context 中提取 Zipkin 的 Span
		traceID := ""
		spanID := ""
		span := zipkin.SpanFromContext(c.Request.Context())

		if span != nil {
			traceID = span.Context().TraceID.String()
			spanID = span.Context().ID.String()
		}

		// 3. 构建日志字段
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("referer", referer),
			zap.Duration("latency", latency),

			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				fields = append(fields, zap.String("gin_error", e))
			}
		}

		// 4. 选择打印类型
		if c.Writer.Status() >= 500 {
			logger.Error("Server Error", fields...)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("Client Error", fields...)
		} else {
			logger.Info("Access Log", fields...)
		}
	}
}
