package kvdb

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Etcd struct {
	cli *clientv3.Client
}

func NewEtcd(ctx context.Context, cfg Config) (h *Etcd, err error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: time.Second * 10,
		TLS:         cfg.tlsConfig(),
		Username:    cfg.Username,
		Password:    cfg.Password,
		Context:     ctx,
	})
	h = &Etcd{cli: client}
	return
}

func (e *Etcd) Get(ctx context.Context, key string) (v []*KV, err error) {
	if strings.HasSuffix(key, "/") {
		return e.GetPrefix(ctx, key)
	}

	resp, err := e.cli.Get(ctx, key)
	if err != nil {
		return
	}

	for _, kv := range resp.Kvs {
		v = append(v, &KV{
			Key:   string(kv.Key),
			Value: g.NewVar(kv.Value),
		})
	}
	return
}

func (e *Etcd) GetPrefix(ctx context.Context, key string) (v []*KV, err error) {
	resp, err := e.cli.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}

	for _, kv := range resp.Kvs {
		v = append(v, &KV{
			Key:   string(kv.Key),
			Value: g.NewVar(kv.Value),
		})
	}
	return
}

func (e *Etcd) Set(ctx context.Context, key string, value interface{}, ttl int64, keepalive ...bool) (err error) {
	opts := make([]clientv3.OpOption, 0)
	if strings.HasSuffix(key, "/") {
		opts = append(opts, clientv3.WithPrefix())
	}
	if ttl > 0 {
		lease := clientv3.NewLease(e.cli)
		var grant *clientv3.LeaseGrantResponse
		if grant, err = lease.Grant(ctx, ttl); err != nil {
			return
		}
		if len(keepalive) > 0 && keepalive[0] {
			go e.keepalive(ctx, lease, grant.ID)
		}
		opts = append(opts, clientv3.WithLease(grant.ID))
	}
	_, err = e.cli.Put(ctx, key, gconv.String(value), opts...)
	return
}

func (e *Etcd) keepalive(ctx context.Context, lease clientv3.Lease, id clientv3.LeaseID) {
	resCh, err := lease.KeepAlive(ctx, id)
	if err != nil {
		return
	}
	for {
		select {
		case _ = <-resCh:
			// discard keepalive message
			// g.Log().Infof(ctx, "Etcd keepalive %v", resp)
		case <-ctx.Done():
			return
		}
	}
}

func (e *Etcd) Delete(ctx context.Context, key string) (err error) {
	opts := make([]clientv3.OpOption, 0)
	if strings.HasSuffix(key, "/") {
		opts = append(opts, clientv3.WithPrefix())
	}
	_, err = e.cli.Delete(ctx, key, opts...)
	return
}

func (e *Etcd) watch(ctx context.Context, key string, handler WatchHandler) {
	opts := make([]clientv3.OpOption, 0)
	if strings.HasSuffix(key, "/") {
		opts = append(opts, clientv3.WithPrefix())
	}
	g.Log().Infof(ctx, "etcd watching %s", key)
	defer func() {
		g.Log().Infof(ctx, "etcd stop watching %s", key)
	}()
	for {
		select {
		case resp := <-e.cli.Watch(ctx, key, opts...):
			if resp.Canceled {
				g.Log().Warningf(ctx, "watch canceled because of: %s", resp.Err())
				return
			}
			for _, ev := range resp.Events {
				var typ EventType
				if ev.IsModify() {
					typ = EventTypeUpdate
				}
				if ev.IsCreate() {
					typ = EventTypeCreate
				}
				if ev.Type == clientv3.EventTypeDelete {
					typ = EventTypeDelete
				}
				handler(ctx, Event{
					KV: KV{
						Key:   string(ev.Kv.Key),
						Value: g.NewVar(ev.Kv.Value),
					},
					Type: typ,
				})
			}
		case <-ctx.Done():
			g.Log().Info(ctx, "etcd watch stop: context canceled")
			return
		}
	}
}

func (e *Etcd) Watch(ctx context.Context, key string, handler WatchHandler) (err error) {
	go e.watch(ctx, key, handler)
	return
}

func (e *Etcd) Locker(ctx context.Context, topic string) (locker sync.Locker, err error) {
	resp, err := e.cli.Grant(ctx, 5)
	if err != nil {
		return
	}
	session, err := concurrency.NewSession(e.cli, concurrency.WithLease(resp.ID))
	if err != nil {
		return
	}
	return concurrency.NewLocker(session, topic), nil
}
