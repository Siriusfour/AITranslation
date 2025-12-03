package types

import (
	"context"
	"github.com/openzipkin/zipkin-go"
)

type TracerInterf interface {
	StartSpanFromContext(ctx context.Context, name string, opts ...zipkin.SpanOption) (zipkin.Span, context.Context)
}
