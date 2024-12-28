package kvdb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func dfs(node *storageNode) {
	if node == nil {
		return
	}
	for _, value := range node.values {
		fmt.Printf("[%s] key=%v value=%v\n", node.name, value.Key, value.Value)
	}
	node.next.Range(func(_, value any) bool {
		dfs(value.(*storageNode))
		return true
	})
}
func TestCachedStorage(t *testing.T) {

	sto := Storages.GetStorage("test")
	cs := sto.(*cachedStorage)
	// print tree
	dfs(cs.root)

	var check = func() bool {
		kvs, err := cs.Get(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		fromDB, err := cs.db.Get(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if len(kvs) != len(fromDB) {
			t.Errorf("key count not match: local=%d, db=%d", len(kvs), len(fromDB))
			return false
		}

		m := make(map[string]*KV)
		for _, kv := range kvs {
			m[kv.Key] = kv
		}

		for _, kv := range fromDB {
			if _, ok := m[kv.Key]; !ok {
				t.Errorf("key not found: %s", kv.Key)
				return false
			}
			if m[kv.Key].Value.String() != kv.Value.String() {
				t.Errorf("value not match: %s , local=%v, db=%v", kv.Key, m[kv.Key].Value, kv.Value)
				return false
			}
			delete(m, kv.Key)
		}

		if len(m) > 0 {
			t.Errorf("local dirty data: %v", m)
			return false
		}
		return true
	}

	fmt.Println("build cache ---------------")
	err := sto.Set(context.Background(), "key", "value")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = sto.Set(context.Background(), "key1", "value1")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = sto.Set(context.Background(), "key2", "value2")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = sto.Set(context.Background(), "key3", "value3")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = sto.Set(context.Background(), "key3/1", "value3-1")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = sto.Set(context.Background(), "key3/2", "value3-2")
	if err != nil {
		t.Fatal(err)
		return
	}
	dfs(cs.root)
	fmt.Println("rebuild cache ---------------")
	cs.buildCache(context.Background())
	dfs(cs.root)
	if !check() {
		t.Fatal("check failed, local cache not equal to db")
	}

	fmt.Println("get all ---------------")
	kvs, err := sto.Get(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, kv := range kvs {
		t.Logf("%+v", kv)
	}
	fmt.Println("get one ---------------")
	vs, err := sto.Get(context.Background(), "key")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("key: %+v", vs[0])

	fmt.Println("delete ---------------")
	err = sto.Delete(context.Background(), "key")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = sto.Delete(context.Background(), "key3/1")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = sto.Delete(context.Background(), "key3")
	if err != nil {
		t.Fatal(err)
		return
	}
	dfs(cs.root)

	if !check() {
		t.Fatal("check failed, local cache not equal to db")
	}
}

func TestEvent(t *testing.T) {
	go func() {
		for {
			dfs(Storages.GetStorage("test").(*cachedStorage).root)
			fmt.Println("---------------------------")
			time.Sleep(time.Second * 5)
		}
	}()

	time.Sleep(time.Minute)
}

func TestStorage_SetTTL(t *testing.T) {
	sto := Storages.GetStorage("test_ttl")
	err := sto.SetTTL(context.Background(), "test_ttl", "value", 10)
	if err != nil {
		t.Fatal(err)
		return
	}

	kvs, err := sto.Get(context.Background(), "test_ttl")
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, kv := range kvs {
		t.Logf("key=%v, value=%v", kv.Key, kv.Value.String())
	}

	if len(kvs) < 1 {
		t.Fatal("get error, no value")
		return
	}

	var checkHasValue = func(kvs []*KV, value string) bool {
		has := false
		for _, kv := range kvs {
			if kv.Value.String() == value {
				has = true
				break
			}
		}
		return has
	}

	if !checkHasValue(kvs, "value") {
		t.Fatal("value not match")
		return
	}

	t.Log("wait 11 seconds for reach ttl")
	time.Sleep(time.Second * 11)

	kvs, err = sto.Get(context.Background(), "test_ttl")
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, kv := range kvs {
		t.Logf("key=%v, value=%v", kv.Key, kv.Value.String())
	}

	if checkHasValue(kvs, "value") {
		t.Fatal("ttl not work")
		return
	}
}
