package structs

import (
	"sync"

	"github.com/gogf/gf/v2/util/gconv"
)

var (
	fieldMappings = sync.Map{} // name:
)

type (
	FieldMapping struct {
		m  map[any]any
		mu sync.RWMutex
	}
)

func SetFieldMapping(name string, m map[any]any) {
	if m == nil {
		return
	}
	value, ok := fieldMappings.Load(name)
	if !ok {
		fieldMappings.Store(name, &FieldMapping{m: m})
		return
	}
	fm := value.(*FieldMapping)
	fm.mu.Lock()
	defer fm.mu.Unlock()
	for k, v := range m {
		fm.m[k] = gconv.PtrAny(v)
	}
	return
}

func GetFieldMappingValue(name string, key any, def ...any) (val any) {
	value, ok := fieldMappings.Load(name)
	if !ok {
		return nil
	}
	fm := value.(*FieldMapping)
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	v, ok := fm.m[key]
	if ok {
		val = v
	} else if len(def) > 0 {
		val = def[0]
	}
	return
}
