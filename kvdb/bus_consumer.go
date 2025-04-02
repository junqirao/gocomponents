package kvdb

import (
	"context"
)

func RegisterHandler(_ context.Context, topic string, handler MessageHandler) {
	handlers.Store(topic, handler)
}
