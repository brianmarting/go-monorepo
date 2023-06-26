package tracing_mocks

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type TracerMock struct {
}

func NewTracerMock() TracerMock {
	return TracerMock{}
}

func (t TracerMock) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx = context.Background()

	return ctx, trace.SpanFromContext(ctx)
}
