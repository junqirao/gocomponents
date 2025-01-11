package objects

import (
	"testing"
)

func TestObjects(t *testing.T) {
	o := O[int]("test")
	o.Set("1", 1)
	o.Set("2", 2)
	o.Range(func(k string, v int) bool {
		t.Log(k, v)
		return true
	})
	res := o.Get("1")
	if res != 1 {
		t.Fatal("expect 1, but got", res)
	}
	res = o.Get("2")
	if res != 2 {
		t.Fatal("expect 2, but got", res)
	}
	o.Delete("1")
	res = o.Get("1")
	if res != 0 {
		t.Fatal("expect 0, but got", res)
	}
	res = o.Get("1", 3)
	if res != 3 {
		t.Fatal("expect 3, but got", res)
	}
}
