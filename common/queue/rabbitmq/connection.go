package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Connection struct {
	channel *amqp.Channel
}

func GetConnection(url string) (Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return Connection{}, err
	}

	channel, err := conn.Channel()
	return Connection{
		channel: channel,
	}, err
}
