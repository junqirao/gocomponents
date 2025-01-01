package structs

import (
	"context"
	"reflect"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	TagHandler func(ctx context.Context, content string, field reflect.StructField, value reflect.Value)
	tagHandler struct {
		name string
		fn   TagHandler
	}
	TagParser struct {
		mu       sync.RWMutex
		handlers []*tagHandler
	}
	TagHandlerOption func(t *TagParser)
)

func NewTagParser(opts ...TagHandlerOption) *TagParser {
	t := &TagParser{}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (p *TagParser) SetHandler(tag string, handler TagHandler) {
	if tag == "" || handler == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	for i, h := range p.handlers {
		if h.name == tag {
			p.handlers[i].fn = handler
			return
		}
	}

	p.handlers = append(p.handlers, &tagHandler{
		name: tag,
		fn:   handler,
	})
}

// Parse struct fields by tag dfs set with SetHandler.
func (p *TagParser) Parse(ctx context.Context, v any) {
	if v == nil {
		return
	}

	var (
		typ = reflect.TypeOf(v)
		val = reflect.ValueOf(v)
	)

	p.parse(ctx, typ, val)
}

func (p *TagParser) parse(ctx context.Context, typ reflect.Type, val reflect.Value) {
	if val.IsNil() {
		return
	}

	if typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Interface {
		typ = typ.Elem()
		val = val.Elem()
	}

	kind := typ.Kind()
	if kind == reflect.Struct || kind == reflect.Ptr {
		for i := 0; i < typ.NumField(); i++ {
			var (
				fieldValue = val.Field(i)
				fieldType  = typ.Field(i)
			)

			p.mu.RLock()
			for _, t := range p.handlers {
				content := fieldType.Tag.Get(t.name)
				if content == "" {
					continue
				}

				t.fn(ctx, content, fieldType, fieldValue)
			}
			p.mu.RUnlock()

			switch fieldValue.Kind() {
			case reflect.Ptr:
				if !fieldValue.IsNil() {
					p.parse(ctx, fieldType.Type, fieldValue)
				}
			case reflect.Struct:
				p.parse(ctx, fieldType.Type, fieldValue)
			default:
			}
		}
	}
}

var (
	WithTagHandlerValueMapping = func() TagHandlerOption {
		return func(t *TagParser) {
			t.SetHandler("mapping", func(ctx context.Context, content string, field reflect.StructField, value reflect.Value) {
				if field.Type.Kind() != reflect.Interface {
					g.Log().Warningf(ctx, "using mapping parser, field %s type is not interface", field.Name)
					return
				}
				if value.IsNil() || !value.CanSet() || value.IsZero() {
					return
				}
				value.Set(reflect.ValueOf(GetFieldMappingValue(content, value.Interface())))
			})
		}
	}
)
