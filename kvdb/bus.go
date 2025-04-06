package kvdb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	busName = "/message-bus"
)

var (
	handlers     = sync.Map{} // topic:MessageHandler
	wgs          = sync.Map{} // id:WaitGroup
	messageCache = sync.Map{} // id:Message
)

func InitBus(ctx context.Context) (err error) {
	pfx := fmt.Sprintf("%s/", busName)
	err = Raw.Watch(ctx, pfx, func(ctx context.Context, e Event) {
		handleMessage(ctx, pfx, e)
	})
	return
}

type (
	MessageHandler interface {
		Handle(ctx context.Context, msg *Message)
		After(ctx context.Context, msg *Message)
	}
	Message struct {
		Id        string `json:"id"`
		Topic     string `json:"topic"`
		Payload   any    `json:"payload"`
		From      string `json:"from"`
		ExpiredAt int64  `json:"expired_at"`
		HasAck    bool   `json:"has_ack"`
		Err       string `json:"err"`
	}
)

func (m Message) Ack(ctx context.Context, err ...error) {
	if m.HasAck {
		return
	}
	if len(err) > 0 && err[0] != nil {
		m.Err = err[0].Error()
		m.HasAck = true
		// deleted at event
		_ = Raw.Set(ctx, buildTopicKey(m.Topic, m.Id), gconv.String(m), WithTTL(10))
	} else {
		_ = Raw.Delete(ctx, buildTopicKey(m.Topic, m.Id))
	}
}

func (m Message) Error() error {
	if m.Err == "" {
		return nil
	}
	return errors.New(m.Err)
}

func buildTopicKey(topic, id string) string {
	return fmt.Sprintf("%s/%s/%s", busName, topic, id)
}

func handleMessage(ctx context.Context, pfx string, e Event) {
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
	case EventTypeCreate:
		msg := &Message{}
		if err := e.Value.Struct(&msg); err != nil {
			g.Log().Warningf(ctx, "invalid message: %s", err.Error())
			return
		}
		if msg.HasAck {
			// drop
			return
		}
		// handle
		handler.Handle(ctx, msg)
	case EventTypeUpdate:
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
		// handle error event or already ack event
		if msg.HasAck || msg.Err != "" {
			_ = Raw.Delete(ctx, buildTopicKey(msg.Topic, msg.Id))
		}
	case EventTypeDelete:
		// finish ack
		mv, ok := messageCache.Load(id)
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
