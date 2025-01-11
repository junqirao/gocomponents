package template

import (
	"testing"
)

func TestTemplate_Parse(t *testing.T) {
	_, err := Create("test", "hello {{.aaa}} {{.Bbb}}")
	if err != nil {
		t.Fatal(err)
		return
	}
	res, err := T("test").Parse(map[string]string{
		"aaa": "123",
		"Bbb": "456",
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("result: ", res)

	if res != "hello 123 456" {
		t.Fatal("unexpected result")
	}
}
