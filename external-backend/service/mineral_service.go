package service

import (
	"common/model"
	"common/observability/tracing"
	"context"
	"external_backend/queue"

	"go.opentelemetry.io/otel/trace"
)

type MineralService interface {
	AddMineral(context.Context, model.Mineral) error
}

type mineralService struct {
	tracer    trace.Tracer
	publisher queue.MineralPublisher
}

func NewMineralService() MineralService {
	return &mineralService{
		tracer:    tracing.GetTracer(),
		publisher: queue.NewMineralPublisher(),
	}
}

func (s mineralService) AddMineral(ctx context.Context, mineral model.Mineral) error {
	spanContext, span := s.tracer.Start(ctx, "mineral-service-add-mineral")
	defer span.End()

	return s.publisher.Publish(spanContext, mineral)
}
