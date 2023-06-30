package queue

type Consumer interface {
	StartConsuming(queueName, key string) (<-chan Message, error)
}
