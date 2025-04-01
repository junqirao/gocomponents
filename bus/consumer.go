package bus

import (
	"context"
	"fmt"
	"sync"

	"github.com/junqirao/gocomponents/kvdb"
)

var (
	handlers     = sync.Map{} // topic:MessageHandler
	wgs          = sync.Map{} // id:WaitGroup
	messageCache = sync.Map{} // id:Message
)

func Init(ctx context.Context) (err error) {
	pfx := fmt.Sprintf("%s/", busName)
	err = kvdb.Raw.Watch(ctx, pfx, func(ctx context.Context, e kvdb.Event) {
		handleMessage(ctx, pfx, e)
	})
	return
}

type (
	MessageHandler interface {
		Handle(ctx context.Context, msg *Message)
		After(ctx context.Context, msg *Message)
	}
)

func RegisterHandler(_ context.Context, topic string, handler MessageHandler) {
	handlers.Store(topic, handler)
}
