package bus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/junqirao/gocomponents/kvdb"
	"github.com/junqirao/gocomponents/registry"
)

const (
	busName = "/message-bus"
)

type Message struct {
	Id        string `json:"id"`
	Topic     string `json:"topic"`
	Payload   any    `json:"payload"`
	From      string `json:"from"`
	ExpiredAt int64  `json:"expired_at"`
	HasAck    bool   `json:"has_ack"`
}

func (m Message) Ack(ctx context.Context) {
	if m.HasAck {
		return
	}
	_ = kvdb.Raw.Delete(ctx, buildTopicKey(m.Topic, m.Id))
}

func buildTopicKey(topic, id string) string {
	return fmt.Sprintf("%s/%s/%s", busName, topic, id)
}

func Push(ctx context.Context, topic string, payload any, ttl int64, wait ...bool) (err error) {
	// check
	_, ok := handlers.Load(topic)
	if !ok {
		err = fmt.Errorf("unknown topic: %s", topic)
		return
	}
	// send
	id := uuid.NewV4().String()
	msg := &Message{
		Id:        id,
		Topic:     topic,
		Payload:   payload,
		From:      registry.Current().Id,
		ExpiredAt: time.Now().Unix() + ttl,
	}
	err = kvdb.Raw.Set(ctx, buildTopicKey(topic, id), gconv.String(msg), ttl)
	if err != nil {
		return
	}
	messageCache.Store(id, msg)
	if !(len(wait) > 0 && wait[0]) {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ttl)*time.Second)
	defer func() {
		wgs.Delete(id)
		cancel()
	}()
	var (
		wg  = &sync.WaitGroup{}
		sig = make(chan struct{})
	)
	go func() {
		wg.Wait()
		close(sig)
	}()

	wg.Add(1)
	wgs.Store(id, wg)
	select {
	case <-ctx.Done():
		err = errors.New("push ack timeout")
	case <-sig:
		// finished
	}
	return
}
