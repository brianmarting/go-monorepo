package queue

import "context"

type Publisher interface {
	Publish(ctx context.Context, key string, data []byte) error
}
