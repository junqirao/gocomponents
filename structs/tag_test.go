package structs

import (
	"context"
	"reflect"
	"testing"
)

type test struct {
	Name  string `tag1:"name" tag2:"xxx"`
	Age   int
	Test1 test3  `tag3:"test1"`
	Test2 test2  `tag3:"test2"`
	Test3 *test2 `tag3:"test3"`
}

type test2 struct {
	Field1 string `tag1:"field1"`
}

type test3 struct {
	Field2 string
	Test4  test2  `tag3:"test4"`
	Test5  *test2 `tag3:"test5"`
}
type test4 struct {
	Field any `mapping:"test"`
}

func TestTagParser(t *testing.T) {
	tp := NewTagParser()
	tp.SetHandler("tag1", func(ctx context.Context, content string, field reflect.StructField, value reflect.Value) {
		t.Logf("tag1: %s, field: %s, value: %v", content, field.Name, value.Interface())
	})
	tp.SetHandler("tag2", func(ctx context.Context, content string, field reflect.StructField, value reflect.Value) {
		t.Logf("tag2: %s, field: %s, value: %v", content, field.Name, value.Interface())
	})
	tp.SetHandler("tag3", func(ctx context.Context, content string, field reflect.StructField, value reflect.Value) {
		t.Logf("tag3: %s, field: %s, value: %v", content, field.Name, value.Interface())
	})
	tp.Parse(context.Background(), test{
		Name: "name",
		Age:  10,
		Test1: test3{
			Field2: "test1",
			Test4: test2{
				Field1: "test4",
			},
		},
		Test2: test2{
			Field1: "test2",
		},
		Test3: &test2{
			Field1: "test3",
		},
	})
}

func TestMappingValue(t *testing.T) {
	SetFieldMapping("test", map[any]any{
		0: "OK",
	})

	t4 := &test4{Field: 0}
	t.Logf("%+v", t4)
	tp := NewTagParser(WithTagHandlerValueMapping())
	tp.Parse(context.Background(), t4)
	t.Logf("%+v", t4)
}
