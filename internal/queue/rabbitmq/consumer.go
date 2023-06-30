package rabbitmq

import (
	"fmt"
	"go-monorepo/internal/queue"

	"github.com/rs/zerolog/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	*Connection
}

func NewConsumer(url string) (queue.Consumer, error) {
	log.Info().Msg("url is " + url)
	conn, err := GetConnection(url)
	if err != nil {
		return consumer{}, fmt.Errorf("failed to get rabbitmq connection: %w", err)
	}

	return consumer{
		Connection: &conn,
	}, nil
}

func (c consumer) StartConsuming(queueName, key string) (<-chan queue.Message, error) {
	_, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = c.channel.QueueBind(
		queueName,
		key,
		"amq.direct",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	amqpDeliveries, err := c.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	messageChannel := make(chan queue.Message)

	go func() {
		defer close(messageChannel)
		for delivery := range amqpDeliveries {
			message := convertToMessage(delivery)
			messageChannel <- message
		}
	}()

	return messageChannel, nil
}

func convertToMessage(delivery amqp.Delivery) Message {
	return Message{
		Delivery: &delivery,
	}
}
