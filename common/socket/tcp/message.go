package tcp

import "common/socket"

type message struct {
	body []byte
}

func NewMessage(body []byte) socket.Message {
	return &message{
		body: body,
	}
}

func (m message) GetBody() []byte {
	return m.body
}
