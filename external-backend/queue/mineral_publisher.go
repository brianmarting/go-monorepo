package queue

import (
	"context"
	"encoding/json"
	facadeQueue "go-monorepo/internal/facade/queue"
	"go-monorepo/internal/model"
	"go-monorepo/internal/observability/tracing"
	"go-monorepo/internal/queue"

	"go.opentelemetry.io/otel/trace"
)

const routingKey = "mineral.deposit"

type MineralPublisher interface {
	Publish(context.Context, model.Mineral) error
}

type mineralPublisher struct {
	tracer    trace.Tracer
	publisher queue.Publisher
}

func NewMineralPublisher() MineralPublisher {
	return &mineralPublisher{
		tracer:    tracing.GetTracer(),
		publisher: facadeQueue.NewPublisher(),
	}
}

func (p *mineralPublisher) Publish(ctx context.Context, model model.Mineral) error {
	spanCtx, span := p.tracer.Start(ctx, "mineral-publisher-marshal-and-publish")
	defer span.End()

	bytes, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return p.publisher.Publish(spanCtx, routingKey, bytes)
}
