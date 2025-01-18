package procedure

import (
	"encoding/json"
	"testing"
)

func TestSetMapValue(t *testing.T) {
	m := map[string]any{
		"key1": "value1",
		"key2": map[string]any{
			"key3": "value3",
		},
	}
	t.Log(SetMapValue(m, "key1", func(v any) any {
		return "new value1"
	}))
	t.Log(SetMapValue(m, "key2.key3", func(v any) any {
		return "new value3"
	}))
	t.Logf("%+v", m)
}

func TestUnmarshal(t *testing.T) {
	m := map[string]any{
		"key1": 1,
		"key2": map[string]any{
			"key3": "3",
		},
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
		return
	}
	res := map[string]any{}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v", res)
}

func BenchmarkSetMapValue(b *testing.B) {
	m := map[string]any{
		"key1": "value1",
		"key2": map[string]any{
			"key3": "value3",
		},
	}
	for i := 0; i < b.N; i++ {
		SetMapValue(m, "key2.key3", "new value3")
	}
}
