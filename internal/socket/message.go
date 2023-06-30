package socket

type Message interface {
	GetBody() []byte
}
