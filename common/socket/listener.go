package socket

type Listener interface {
	Start(done <-chan interface{}) (<-chan Message, chan<- string)
}
