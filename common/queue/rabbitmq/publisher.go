package rabbitmq

import (
	"common/queue"
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type publisher struct {
	*Connection
}

func NewPublisher(url string) (queue.Publisher, error) {
	conn, err := GetConnection(url)
	if err != nil {
		return publisher{}, nil
	}

	return publisher{
		Connection: &conn,
	}, nil
}

func (p publisher) Publish(ctx context.Context, routingKey string, data []byte) error {
	return p.channel.PublishWithContext(
		ctx,
		"amq.direct",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Timestamp:   time.Now(),
		},
	)
}
