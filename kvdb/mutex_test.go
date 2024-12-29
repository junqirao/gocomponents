package kvdb

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	ctx := context.Background()
	mutex, err := NewMutex(ctx, "test")
	if err != nil {
		t.Fatal(err)
		return
	}
	mutex.Lock()
	t.Log("lock 1")
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		t.Log("try lock 2")
		mutex.Lock()
		time.Sleep(time.Second * 1)
		t.Log("lock 2")
		mutex.Unlock()
		t.Log("unlock 2")
		wg.Done()
	}()
	time.Sleep(time.Second * 5)
	mutex.Unlock()
	t.Log("unlock 1")
	wg.Wait()
}
