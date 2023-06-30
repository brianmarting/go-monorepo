package rabbitmq

import (
	"context"
	"go-monorepo/internal/observability/tracing"
	"go-monorepo/internal/queue"
	"time"

	"go.opentelemetry.io/otel/trace"

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
	span := trace.SpanFromContext(ctx)

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
			Headers: map[string]interface{}{
				tracing.TraceIdHeader: span.SpanContext().TraceID().String(),
				tracing.SpanIdHeader:  span.SpanContext().SpanID().String(),
			},
		},
	)
}
