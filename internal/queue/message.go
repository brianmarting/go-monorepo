package queue

type Message interface {
	GetBytes() []byte
	GetHeaders() map[string]interface{}
	Ack() error
	Nack() error
}
