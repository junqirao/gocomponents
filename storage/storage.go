package storage

import (
	"context"
	"io"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	defaultStorageName = "default"
)

type (
	IStorage interface {
		Put(ctx context.Context, name string, r io.Reader) (key string, err error)
		Get(ctx context.Context, name string) (rc io.ReadCloser, err error)
		Delete(ctx context.Context, name string) (err error)
		SignGetUrl(ctx context.Context, name string, expires int64, contentType string, disposition string) (s string, err error)
		SignPutUrl(ctx context.Context, name string, expires int64) (s string, err error)
	}
)

var (
	ins = sync.Map{} // name : IStorage
)

func Init(ctx context.Context) (err error) {
	v, err := g.Cfg().Get(ctx, "storage")
	if err != nil {
		return
	}
	m := make(map[string]Config)
	if err = v.Struct(&m); err != nil {
		return
	}

	for name, cfg := range m {
		var storage IStorage
		if storage, err = New(cfg); err != nil {
			return
		}
		ins.Store(name, storage)
	}
	return
}

func MustInit(ctx context.Context) {
	if err := Init(ctx); err != nil {
		panic(err)
	}
}

func Storage(name ...string) IStorage {
	n := defaultStorageName
	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}
	v, ok := ins.Load(n)
	if !ok {
		return insEmpty
	}
	return v.(IStorage)
}

func Default() IStorage {
	return Storage(defaultStorageName)
}

func New(config Config) (s IStorage, err error) {
	switch config.Type {
	case TypeMinio:
		s, err = newMinio(config)
	default:
		// using minio
		s, err = newMinio(config)
	}
	return
}
