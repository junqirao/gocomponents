package bus

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"

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
		key := strings.TrimPrefix(e.Key, pfx)
		parts := strings.Split(key, "/")
		if len(parts) != 2 {
			g.Log().Warningf(ctx, "invalid topic: %s", key)
			return
		}
		var (
			topic = parts[0]
			id    = parts[1]
		)

		v, ok := handlers.Load(topic)
		if !ok {
			// drop
			return
		}
		handler := v.(MessageHandler)
		switch e.Type {
		case kvdb.EventTypeCreate:
			msg := &Message{}
			if err := e.Value.Struct(&msg); err != nil {
				g.Log().Warningf(ctx, "invalid message: %s", err.Error())
				return
			}
			// start
			err := handler.Receive(ctx, msg)
			if err != nil {
				// drop
				return
			}
			msg.Ack(ctx)
		case kvdb.EventTypeDelete:
			// finish ack
			mv, ok := messageCache.LoadAndDelete(id)
			if !ok {
				// drop
				return
			}
			if wgv, loaded := wgs.LoadAndDelete(id); loaded {
				wgv.(*sync.WaitGroup).Done()
			}
			handler.Finish(ctx, mv.(*Message))
		}
	})
	return
}

type (
	MessageHandler interface {
		Receive(ctx context.Context, msg *Message) (err error)
		Finish(ctx context.Context, msg *Message)
	}
)

func RegisterHandler(_ context.Context, topic string, handler MessageHandler) {
	handlers.Store(topic, handler)
}
