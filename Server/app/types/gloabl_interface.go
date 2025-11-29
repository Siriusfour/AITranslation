package types

import (
	"context"
	"github.com/openzipkin/zipkin-go"
)

type TracerInterf interface {
	StartSpanFromContext(context.Context, string) (zipkin.Span, context.Context)
}
