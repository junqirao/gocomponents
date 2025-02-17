package kvdb

import (
	"context"
	"strings"
)

type (
	storage struct {
		prefix string
		cfg    StorageConfig
		name   string
		Database
	}
)

func newStorage(prefix, name string, db Database, cfg StorageConfig) *storage {
	name = strings.ReplaceAll(name, cfg.Separator, "")
	return &storage{
		prefix:   prefix,
		cfg:      cfg,
		name:     name,
		Database: db,
	}
}

func (s *storage) Get(ctx context.Context, key ...string) (v []*KV, err error) {
	return s.Database.GetPrefix(ctx, s.buildStorageKey(key...))
}

func (s *storage) Set(ctx context.Context, key string, value interface{}) (err error) {
	return s.set(ctx, key, value, 0)
}

func (s *storage) SetTTL(ctx context.Context, key string, value interface{}, ttl int64, keepalive ...bool) (err error) {
	return s.set(ctx, key, value, ttl, keepalive...)
}

func (s *storage) set(ctx context.Context, key string, value interface{}, ttl int64, keepalive ...bool) (err error) {
	if !strings.HasPrefix(key, s.name) {
		key = s.buildStorageKey(key)
	}
	return s.Database.Set(ctx, key, value, ttl, keepalive...)
}

func (s *storage) Delete(ctx context.Context, key string) (err error) {
	return s.Database.Delete(ctx, s.buildStorageKey(key))
}

func (s *storage) buildStorageKey(key ...string) string {
	builder := strings.Builder{}
	builder.WriteString(s.prefix)
	builder.WriteString(s.name)
	builder.WriteString(s.cfg.Separator)
	if len(key) > 0 && key[0] != "" {
		builder.WriteString(key[0])
	}
	return builder.String()
}
