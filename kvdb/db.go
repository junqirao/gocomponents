package kvdb

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	Raw              Database
	databaseOnceInit sync.Once
	storageOnceInit  sync.Once
)

type (
	// WatchHandler watch registry
	WatchHandler func(ctx context.Context, e Event)
	// KV kv
	KV struct {
		Key   string
		Value *g.Var
	}
	// Event of database key value changes
	Event struct {
		KV
		Type EventType
	}
	// Database abstract key-value database ability
	Database interface {
		// Get values from database by key
		Get(ctx context.Context, key string) (v []*KV, err error)
		// GetPrefix values from database by prefixed key
		GetPrefix(ctx context.Context, key string) (v []*KV, err error)
		// Set value to database
		Set(ctx context.Context, key string, value interface{}, opts ...SetOption) (err error)
		// Delete value from database
		Delete(ctx context.Context, key string) (err error)
		// Watch database changes
		Watch(ctx context.Context, key string, handler WatchHandler) (err error)
		// Locker distribute lock
		Locker(ctx context.Context, topic string) (locker sync.Locker, err error)
	}
	// EventType of instance change
	EventType string

	// SetOption set handler
	SetOption func(cfg *setConfig)
	// keepAliveHandler
	setConfig struct {
		keepalive          bool
		ttl                int64
		onKeepaliveStopped func(err error)
	}
)

var (
	WithKeepAlive = func() SetOption {
		return func(cfg *setConfig) {
			cfg.keepalive = true
		}
	}
	WithTTL = func(ttlSeconds int64) SetOption {
		return func(cfg *setConfig) {
			cfg.ttl = ttlSeconds
		}
	}
	WithKeepAliveStoppedHandler = func(f func(err error)) SetOption {
		return func(cfg *setConfig) {
			cfg.onKeepaliveStopped = f
		}
	}
)

// event type define
const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "upsert"
	EventTypeDelete EventType = "delete"
)
