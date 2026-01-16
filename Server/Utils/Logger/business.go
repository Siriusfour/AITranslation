package Logger

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	"go.uber.org/zap"
)

// GetLoggerWithTrace 从 map 中获取指定 name 的 logger，并注入 trace_id
func GetLoggerWithTrace(ctx context.Context, logMap map[string]*zap.Logger, name string) *zap.Logger {
	// 1. 从 map 获取基础 logger
	l, ok := logMap[name]
	if !ok {
		// 兜底，防止 map key 写错导致 panic
		l = zap.L() // 或者 zap.NewNop()
	}

	// 2. 从 context 获取 Zipkin Span
	span := zipkin.SpanFromContext(ctx)
	if span == nil {
		// 如果当前没有 trace (比如后台定时任务)，直接返回原始 logger
		return l
	}

	// 3. 注入 trace_id 和 span_id
	// zap.With 是轻量级操作（Copy-on-write），性能很高
	return l.With(
		zap.String("trace_id", span.Context().TraceID.String()),
		zap.String("span_id", span.Context().ID.String()),
	)
}
