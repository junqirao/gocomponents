package kvdb

import (
	"context"
)

func RegisterBusHandler(_ context.Context, topic string, handler MessageHandler) {
	handlers.Store(topic, handler)
}
