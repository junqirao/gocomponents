package procedure

import (
	"github.com/gogf/gf/v2/container/gvar"
)

type (
	Meta map[string]any
)

func (m Meta) Get(key string, def ...any) (val *gvar.Var) {
	val = gvar.New(nil)
	if len(def) > 0 && def[0] != nil {
		val.Set(def[0])
	}
	if m == nil {
		return
	}
	if v, ok := m[key]; ok {
		val.Set(v)
	}
	return
}
