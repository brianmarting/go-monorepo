package queue

import (
	"fmt"
	"go-monorepo/internal/queue"
	"go-monorepo/internal/queue/rabbitmq"
	"os"

	"github.com/rs/zerolog/log"
)

func NewPublisher() queue.Publisher {
	pub, err := rabbitmq.NewPublisher(getUrl())
	if err != nil {
		log.Panic().Msg("failed to create publisher")
	}

	return pub
}

func NewConsumer() queue.Consumer {
	consumer, err := rabbitmq.NewConsumer(getUrl())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create consumer")
	}

	return consumer
}

func getUrl() string {
	var (
		username = os.Getenv("RABBITMQ_USERNAME")
		password = os.Getenv("RABBITMQ_PASSWORD")
		host     = os.Getenv("RABBITMQ_HOST")
	)

	return fmt.Sprintf("amqp://%s:%s@%s:5672/", username, password, host)
}
