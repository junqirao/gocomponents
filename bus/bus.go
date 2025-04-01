package bus

import (
	"context"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/kvdb"
)

func handleMessage(ctx context.Context, pfx string, e kvdb.Event) {
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
		err := handler.Handle(ctx, msg)
		if err != nil {
			// drop
			return
		}
	case kvdb.EventTypeUpdate:
		// extra
		_, ok = messageCache.Load(id)
		if !ok {
			return
		}
		msg := &Message{}
		if err := e.Value.Struct(msg); err != nil {
			// drop
			g.Log().Warningf(ctx, "invalid message: %s", err.Error())
			return
		}
		messageCache.Store(id, msg)
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
		msg := mv.(*Message)
		msg.HasAck = true
		handler.After(ctx, msg)
	}
}
