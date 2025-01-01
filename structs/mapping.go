package structs

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
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

type (
	mappingStorage     []mappingStorageUnit
	mappingStorageUnit struct {
		Name    string                   `json:"name"`
		Content []mappingStorageKeyValue `json:"content"`
	}
	mappingStorageKeyValue struct {
		Key       any       `json:"key"`
		KeyType   fieldType `json:"key_type"`
		Val       any       `json:"value"`
		ValueType fieldType `json:"value_type"`
	}
	fieldType string
)

const (
	fieldTypeString fieldType = "string"
	fieldTypeInt    fieldType = "int"
	fieldTypeUint   fieldType = "uint"
	fieldTypeFloat  fieldType = "float"
)

func (u mappingStorageUnit) buildMap() map[any]any {
	m := map[any]any{}
	for _, kv := range u.Content {
		m[kv.GetKey()] = kv.GetValue()
	}
	return m
}

func (k mappingStorageKeyValue) GetKey() any {
	return k.parse(k.KeyType, k.Key)
}

func (k mappingStorageKeyValue) GetValue() any {
	return k.parse(k.ValueType, k.Val)
}

func (k mappingStorageKeyValue) parse(typ fieldType, val any) any {
	var vv = val
	switch typ {
	case fieldTypeString:
		vv = gconv.String(val)
	case fieldTypeInt:
		vv = gconv.Int(val)
	case fieldTypeUint:
		vv = gconv.Uint(val)
	case fieldTypeFloat:
		vv = gconv.Float64(val)
	}
	return vv
}

func LoadMappingFromEmbed(ctx context.Context, efs embed.FS) (err error) {
	dir, err := efs.ReadDir("mapping")
	if err != nil {
		return
	}
	for _, entry := range dir {
		if !entry.IsDir() &&
			!strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		var (
			info    fs.FileInfo
			content []byte
		)

		if info, err = entry.Info(); err != nil {
			return
		}
		if content, err = efs.ReadFile("mapping/" + info.Name()); err != nil {
			return
		}
		storage := mappingStorage{}
		if err = json.Unmarshal(content, &storage); err != nil {
			return
		}

		for _, unit := range storage {
			m := unit.buildMap()
			SetFieldMapping(unit.Name, m)
			g.Log().Infof(ctx, "structs field mapping loaded from: %s, %s (%v)", info.Name(), unit.Name, len(m))
		}
	}
	return
}
