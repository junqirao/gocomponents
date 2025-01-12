package structs

import (
	"context"
	"reflect"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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

func (p *TagParser) TryParse(ctx context.Context, v any) {
	if v == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			g.Log().Warningf(ctx, "tag parser try parse error: %v", r)
		}
	}()

	p.Parse(ctx, v)
}

func (p *TagParser) parse(ctx context.Context, typ reflect.Type, val reflect.Value) {
	if typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Interface {
		if val.IsNil() {
			return
		}

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
				if value.IsNil() || !value.CanSet() {
					return
				}
				v := GetFieldMappingValue(content, value.Interface())
				if v == nil {
					return
				}
				value.Set(reflect.ValueOf(v))
			})
		}
	}
	// WithTagHandlerDefaultVal set default value for field when value not zero,
	// supports type of string, int8/16/32/64, uint8/16/32/64, float32/64, bool.
	WithTagHandlerDefaultVal = func() TagHandlerOption {
		return func(t *TagParser) {
			t.SetHandler("default", func(ctx context.Context, content string, field reflect.StructField, value reflect.Value) {
				kind := field.Type.Kind()
				if !value.CanSet() || !value.IsZero() {
					return
				}
				switch kind {
				case reflect.String:
					value.SetString(content)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value.SetInt(gconv.Int64(content))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					value.SetUint(gconv.Uint64(content))
				case reflect.Float32, reflect.Float64:
					value.SetFloat(gconv.Float64(content))
				case reflect.Bool:
					value.SetBool(gconv.Bool(content))
				default:
				}
			})
		}
	}
)
