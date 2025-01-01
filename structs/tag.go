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

	// if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
	// 	val = val.Elem()
	// } else {
	// 	val = reflect.ValueOf(v)
	// }
	//
	// fmt.Printf("val: isvalid:%v canaddr:%v\n", val.IsValid(), val.CanAddr())
	//
	// if kind := typ.Kind(); kind == reflect.Struct {
	// 	for i := 0; i < typ.NumField(); i++ {
	// 		var (
	// 			value = val.Field(i)
	// 			field = typ.Field(i)
	// 		)
	//
	// 		p.mu.RLock()
	// 		for _, t := range p.handlers {
	// 			content := field.Tag.Get(t.name)
	// 			if content == "" {
	// 				continue
	// 			}
	//
	// 			fmt.Printf("isvalid:%v canaddr:%v\n", value.IsValid(), value.CanAddr())
	// 			t.fn(ctx, content, field, value)
	// 		}
	// 		p.mu.RUnlock()
	//
	// 		switch field.Type.Kind() {
	// 		case reflect.Ptr:
	// 			if !value.IsNil() {
	// 				p.Parse(ctx, value.Elem().Interface())
	// 			}
	// 		case reflect.Struct:
	// 			p.Parse(ctx, value.Interface())
	// 		default:
	// 		}
	//
	// 	}.

	// } else if kind == reflect.Ptr {
	// 	p.Parse(ctx, val.Elem().Interface())
	// }
}

func (p *TagParser) parse(ctx context.Context, typ reflect.Type, val reflect.Value) {
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
				if value.IsNil() {
					return
				}
				value.Set(reflect.ValueOf(GetFieldMappingValue(content, value.Interface())))
			})
		}
	}
)
