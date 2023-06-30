package tcp

import (
	"bufio"
	"fmt"
	"go-monorepo/internal/socket"
	"net"

	"github.com/rs/zerolog/log"
)

type listener struct {
	port string
}

func NewListener(port string) socket.Listener {
	return listener{
		port: port,
	}
}

func (l listener) Start(done <-chan interface{}) (<-chan socket.Message, chan<- string) {
	in := make(chan socket.Message)
	out := make(chan string)

	go start(l.port, in, out, done)

	return in, out
}

func start(port string, in chan<- socket.Message, out <-chan string, done <-chan interface{}) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Error().Err(err).Msg("failed to start tcp listener")
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept listener")
			return
		}

		log.Info().Msg(fmt.Sprintf("started listening on tcp port %s", port))

		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
				msg, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					log.Error().Err(err).Msg("something went wrong when reading connection, disconnecting...")
					return
				}

				in <- convertToMessage(msg)

				resp := <-out
				if _, err := conn.Write([]byte(resp)); err != nil {
					log.Info().Err(err).Msg("something went wrong when sending tcp msg reply")
				}
			}
		}()
	}
}

func convertToMessage(message string) socket.Message {
	return NewMessage(
		[]byte(message),
	)
}
