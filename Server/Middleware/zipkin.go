package Middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go" // 引入 zipkin 官方库
	"go.uber.org/zap"
)

// HttpLog  来记录 HTTP Access Log
func HttpLog(logMap map[string]*zap.Logger, logKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 map 中安全获取指定的 logger
		// 如果 key 不存在，可以使用默认的 no-op logger 防止 panic，或者直接 panic 提醒配置错误
		logger, ok := logMap[logKey]
		if !ok {
			// 兜底方案：如果没有找到对应的 logger，创建一个临时的，避免空指针崩溃
			// 实际生产中建议在启动时就检查 map 完整性
			logger = zap.NewNop()
		}

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		// --- 请求结束，开始记录日志 ---

		end := time.Now()
		latency := end.Sub(start)

		// 2. 【核心修改】适配 Zipkin
		// 从 Context 中提取 Zipkin 的 Span
		traceID := ""
		spanID := ""
		span := zipkin.SpanFromContext(c.Request.Context())

		// 必须判断 span 是否为 nil，因为有些请求可能没开启追踪
		if span != nil {
			// Zipkin 的 TraceID 获取方式
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
			zap.Duration("latency", latency),
			// 注入 Zipkin 的 ID
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
			// 标记这是 HTTP Access Log
			zap.String("log_type", "access"),
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
