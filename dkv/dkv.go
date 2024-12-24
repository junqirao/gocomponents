package dkv

import (
	"context"
	"errors"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	WatchHandler func(ctx context.Context, e Event)
	KV           struct {
		Key   string
		Value *g.Var
	}
	Event struct {
		KV
		Type EventType
	}
	EventType string
	DBType    string
)

const (
	TypeETD DBType = "etcd"
)

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

type Database interface {
	Get(ctx context.Context, key string) (v []*KV, err error)
	GetPrefix(ctx context.Context, key string) (v []*KV, err error)
	Set(ctx context.Context, key string, value interface{}, ttl ...int64) (err error)
	Delete(ctx context.Context, key string) (err error)
	Watch(ctx context.Context, key string, handler WatchHandler) (err error)
	Locker(ctx context.Context, topic string) (locker sync.Locker, err error)
}

func NewDB(ctx context.Context, typ ...DBType) (db Database, err error) {
	cfg := new(Config)
	val, err := g.Cfg().Get(ctx, "registry")
	if err != nil {
		g.Log().Warningf(ctx, "failed to get registry config: %v", err)
		return
	}
	if err = val.Struct(&cfg); err != nil {
		g.Log().Warningf(ctx, "failed to parse registry config: %v", err)
		return
	}

	return NewDBWithConfig(ctx, cfg, typ...)
}

func NewDBWithConfig(ctx context.Context, cfg *Config, typ ...DBType) (db Database, err error) {
	typeName := TypeETD
	if len(typ) > 0 && typ[0] != "" {
		typeName = typ[0]
	}

	switch typeName {
	case TypeETD:
		db, err = newEtcd(ctx, *cfg)
	default:
		err = errors.New("unsupported database type: " + string(typeName))
	}
	return
}
