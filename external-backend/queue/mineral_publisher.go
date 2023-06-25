package queue

import (
	facadeQueue "common/facade/queue"
	"common/model"
	"common/queue"
	"context"
	"encoding/json"
)

const routingKey = "mineral.deposit"

type MineralPublisher interface {
	Publish(context.Context, model.Mineral) error
}

type mineralPublisher struct {
	publisher queue.Publisher
}

func NewMineralPublisher() MineralPublisher {
	return &mineralPublisher{
		publisher: facadeQueue.NewPublisher(),
	}
}

func (p *mineralPublisher) Publish(ctx context.Context, model model.Mineral) error {
	bytes, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, routingKey, bytes)
}
