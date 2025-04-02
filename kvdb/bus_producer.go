package kvdb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/junqirao/gocomponents/meta"
)

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
		From:      meta.InstanceId(),
		ExpiredAt: time.Now().Unix() + ttl,
	}
	err = Raw.Set(ctx, buildTopicKey(topic, id), gconv.String(msg), ttl)
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
		value, loaded := messageCache.LoadAndDelete(id)
		if loaded {
			msg = value.(*Message)
			err = msg.Error()
		}
	}
	return
}
