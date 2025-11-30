package Middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go" // 引入 zipkin 包
	"go.uber.org/zap"
)

// GinAccessLog 接收你的 logger map
func GinAccessLog(loggers map[string]*zap.Logger) gin.HandlerFunc {
	// 1. 安全获取需要的 logger 实例
	// 假设 HTTP 请求日志使用 map 中的 "Business" 实例
	accessLogger, ok := loggers["Business"]
	if !ok {
		// 降级策略：如果没有找到，使用 Nop (什么都不干) 或者 zap.L()
		// 建议你在 main 初始化时保证 key 存在
		accessLogger = zap.NewNop()
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 执行后续逻辑
		c.Next()

		// --- 请求处理完毕 ---
		end := time.Now()
		latency := end.Sub(start)

		// 2. 【关键】提取 Zipkin TraceID
		traceID := ""
		spanID := ""
		span := zipkin.SpanFromContext(c.Request.Context())
		if span != nil {
			// Zipkin 的 SpanContext 包含 TraceID
			ctx := span.Context()
			traceID = ctx.TraceID.String()
			spanID = ctx.ID.String()
		}

		// 3. 构建日志字段
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			// 注入 TraceID，方便在 Kibana 里搜索
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
		}

		// 4. 处理 Gin 内部错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				fields = append(fields, zap.String("gin_error", e))
			}
		}

		// 5. 分级打印
		if c.Writer.Status() >= 500 {
			accessLogger.Error("Server Error", fields...)
		} else if c.Writer.Status() >= 400 {
			accessLogger.Warn("Client Error", fields...)
		} else {
			accessLogger.Info("Access Log", fields...)
		}
	}
}
