package kvdb

import (
	"context"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	// Storages Global instance of storages
	Storages *storages
)

type (
	// Storage interface
	Storage interface {
		// Get value
		Get(ctx context.Context, key ...string) (v []*KV, err error)
		// Set value
		Set(ctx context.Context, key string, value interface{}) (err error)
		// SetTTL set value with ttl in second
		SetTTL(ctx context.Context, key string, value interface{}, ttl int64, keepalive ...bool) (err error)
		// Delete value
		Delete(ctx context.Context, key string) (err error)
	}
	// StorageEventHandler process storage event
	StorageEventHandler func(t EventType, key string, value interface{})
	storages            struct {
		ctx context.Context
		cfg Config
		db  Database
		m   sync.Map // key: (name)string, value: Storage
		evs sync.Map // key: (name)string, value: StorageEventHandler
	}
)

func newStorages(ctx context.Context, cfg Config, db Database) *storages {
	sto := &storages{ctx: ctx, cfg: cfg, db: db}
	// watch and update caches event bus
	sto.watchAndUpdateCaches(ctx)
	return sto
}

// GetStorage or create Storage instance
func (s *storages) GetStorage(name string, uncached ...bool) Storage {
	var cs *cachedStorage
	v, ok := s.m.Load(name)
	if ok {
		cs = v.(*cachedStorage)
	}

	if cs == nil {
		cs = newCachedStorage(s.ctx, newStorage(name, s.db, s.cfg))
		s.m.Store(name, cs)
	}

	if len(uncached) > 0 && uncached[0] {
		return cs.db
	}
	return cs
}

func (s *storages) watchAndUpdateCaches(ctx context.Context) {
	pfx := s.cfg.Prefix
	err := s.db.Watch(ctx, pfx, func(ctx context.Context, e Event) {
		pos := strings.Split(strings.TrimPrefix(e.Key, pfx), s.cfg.Separator)
		if len(pos) == 0 {
			return
		}
		name := pos[0]
		key := strings.Join(pos[1:], s.cfg.Separator)

		// internal event
		if sto, ok := s.m.Load(name); ok {
			sto.(*cachedStorage).handleEvent(e.Type, key, e.Value)
		}

		// push to event handler
		if ev, ok := s.evs.Load(name); ok {
			ev.(StorageEventHandler)(e.Type, key, e.Value)
		}
	})
	if err != nil {
		g.Log().Errorf(ctx, "failed to watch and update caches at storage: %s", err.Error())
	}
}

func (s *storages) SetEventHandler(name string, handler StorageEventHandler) {
	s.evs.Store(name, handler)
}
