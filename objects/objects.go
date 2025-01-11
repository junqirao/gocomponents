package objects

import (
	"sync"
)

var (
	instances = sync.Map{} // name:*Objects
)

func O[T any](name string) *Objects[T] {
	v, ok := instances.Load(name)
	if ok {
		return v.(*Objects[T])
	}
	ins := &Objects[T]{name: name}
	instances.Store(name, ins)
	return ins
}

type Objects[T any] struct {
	name string
	m    sync.Map
}

func (o *Objects[T]) Get(k string, def ...T) (res T) {
	v, ok := o.m.Load(k)
	if ok {
		if res, ok = v.(T); ok {
			return
		}
	}

	if len(def) > 0 {
		res = def[0]
	}

	return
}

func (o *Objects[T]) Set(k string, v T) {
	o.m.Store(k, v)
}

func (o *Objects[T]) Delete(k string) {
	o.m.Delete(k)
}

func (o *Objects[T]) Range(f func(k string, v T) bool) {
	o.m.Range(func(k, v any) bool {
		return f(k.(string), v.(T))
	})
}
