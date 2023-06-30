package queue

import (
	facadeQueue "go-monorepo/internal/facade/queue"
	"go-monorepo/internal/queue"
)

type MineralConsumer interface {
	StartConsuming(queueName, key string) (<-chan queue.Message, error)
}

type mineralConsumer struct {
	consumer queue.Consumer
}

func NewMineralConsumer() MineralConsumer {
	return &mineralConsumer{
		consumer: facadeQueue.NewConsumer(),
	}
}

func (c mineralConsumer) StartConsuming(queueName, key string) (<-chan queue.Message, error) {
	return c.consumer.StartConsuming(queueName, key)
}
