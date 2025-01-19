package procedure

import (
	"context"
	"strings"
	"text/template"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	FuncNameSetInput    = "set_input"
	FuncNameGetInput    = "get_input"
	FuncNameSetResult   = "set_result"
	FuncNameGetResult   = "get_result"
	FuncNameNewMap      = "new_map"
	FuncNameSetMapValue = "set_map_value"
	FuncNameLogInfo     = "info"
	FuncNameLogWarning  = "warning"
	FuncNameLogError    = "error"
)

func setNodeDefaultFunctions(ctx context.Context, node *Node, fm template.FuncMap, input *gmap.StrAnyMap, results *gmap.Map, out any) {
	fm[FuncNameSetInput] = func(k string, v any) error {
		input.Set(k, v)
		return nil
	}
	fm[FuncNameGetInput] = func(k string) any {
		return input.Get(k)
	}
	fm[FuncNameSetResult] = func(k string, v any) error {
		if k == "" {
			k = node.Name
		}
		results.Set(k, v)
		return nil
	}
	fm[FuncNameGetResult] = func(k string) any {
		return results.Get(k)
	}
	setDefaultFunctions(ctx, fm)
}

func setDefaultFunctions(ctx context.Context, fm template.FuncMap) {
	fm[FuncNameLogInfo] = func(format string, v ...any) any {
		glog.Infof(ctx, format, v...)
		return nil
	}
	fm[FuncNameLogWarning] = func(format string, v ...any) any {
		glog.Warningf(ctx, format, v...)
		return nil
	}
	fm[FuncNameLogError] = func(format string, v ...any) any {
		glog.Errorf(ctx, format, v...)
		return nil
	}
	fm[FuncNameSetMapValue] = SetMapValue
	fm[FuncNameNewMap] = func() map[string]any {
		return make(map[string]any)
	}

}

func SetMapValue(m any, index string, v any) any {
	if m == nil || index == "" {
		return nil
	}
	return setMapValue(m, v, strings.Split(index, ".")...)
}

func setMapValue(mv any, v any, indexes ...string) any {
	if mv == nil || len(indexes) == 0 {
		return nil
	}

	m := gconv.Map(mv)
	if len(indexes) == 1 {
		m[indexes[0]] = v
		return v
	}
	for i, index := range indexes {
		val, ok := m[index]
		if !ok {
			return nil
		}
		if i == len(indexes)-1 {
			m[index] = v
			return v
		}
		m = gconv.Map(val)
	}

	return nil
}
