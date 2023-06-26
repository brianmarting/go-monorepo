package service

import (
	"common/model"
	"common/observability/tracing"
	"context"
	"external-backend/queue"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

var (
	mineralServiceOnce     sync.Once
	mineralServiceInstance *mineralService
)

type MineralService interface {
	AddMineral(context.Context, model.Mineral) error
}

type mineralService struct {
	tracer    trace.Tracer
	publisher queue.MineralPublisher
}

func GetMineralService() MineralService {
	mineralServiceOnce.Do(func() {
		mineralServiceInstance = &mineralService{
			tracer:    tracing.GetTracer(),
			publisher: queue.NewMineralPublisher(),
		}
	})
	return mineralServiceInstance
}

func (s mineralService) AddMineral(ctx context.Context, mineral model.Mineral) error {
	spanContext, span := s.tracer.Start(ctx, "mineral-service-add-mineral")
	defer span.End()

	return s.publisher.Publish(spanContext, mineral)
}
