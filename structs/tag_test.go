package structs

import (
	"context"
	"embed"
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

type test5 struct {
	Field1 int8    `default:"127"`
	Field2 uint16  `default:"65535"`
	Field3 float64 `default:"3.14"`
	Field4 bool    `default:"true"`
	Field5 string  `default:"hello,world"`
	Field6 int8    `default:"128"`
	Field7 int8    `default:"100"`
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
		1: 0,
	})

	t4 := &test4{Field: 0}
	t.Logf("%+v", t4)
	tp := NewTagParser(WithTagHandlerValueMapping())
	tp.Parse(context.Background(), t4)
	t.Logf("%+v", t4)
	t.Log("==========================")
	var t5 any = &test4{Field: 0}
	t.Logf("%+v", t5)
	tp.Parse(context.Background(), t5)
	t.Logf("%+v", t5)
	t.Log("==========================")
	t5 = any(&test4{Field: 0})
	t.Logf("%+v", t5)
	tp.Parse(context.Background(), &t5)
	t.Logf("%+v", t5)
	t.Log("==========================")
	t5 = reflect.ValueOf(any(&test4{Field: 0})).Interface()
	t.Logf("%+v", t5)
	tp.Parse(context.Background(), &t5)
	t.Logf("%+v", t5)
	t.Log("==========================")
	t5 = reflect.ValueOf(any(&test4{Field: 0})).Interface()
	t.Logf("%+v", t5)
	tp.Parse(context.Background(), t5)
	t.Logf("%+v", t5)
	t.Log("==========================")
	t5 = any(&test4{Field: 1})
	t.Logf("%+v", t5)
	tp.Parse(context.Background(), t5)
	t.Logf("%+v", t5)
}

//go:embed mapping
var efs embed.FS

func TestLoadMappingFromEmbed(t *testing.T) {
	err := LoadMappingFromEmbed(context.Background(), efs)
	if err != nil {
		t.Fatal(err)
	}
	val := GetFieldMappingValue("test0", 0)
	t.Logf("%+v", val)
}

func TestDefaultValue(t *testing.T) {
	parser := NewTagParser(WithTagHandlerDefaultVal())
	t5 := &test5{
		Field7: 99,
	}
	parser.Parse(context.Background(), t5)
	t.Logf("%+v", t5)
	if t5.Field1 != 127 {
		t.Fatal("field1 not equal")
	}
	if t5.Field2 != 65535 {
		t.Fatal("field2 not equal")
	}
	if t5.Field3 != 3.14 {
		t.Fatal("field3 not equal")
	}
	if t5.Field4 != true {
		t.Fatal("field4 not equal")
	}
	if t5.Field5 != "hello,world" {
		t.Fatal("field5 not equal")
	}
	if t5.Field6 != -128 {
		t.Fatal("field6 not equal")
	}
	if t5.Field7 != 99 {
		t.Fatal("field7 not equal")
	}
}
