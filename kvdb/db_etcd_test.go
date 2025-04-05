package kvdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func TestKeepAlive(t *testing.T) {
	go func() {
		s := g.Server()
		s.EnablePProf()
		s.SetPort(9999)
		s.Run()
	}()
	cli, err := NewEtcd(context.Background(), Config{
		Endpoints: []string{"172.18.28.241:2379", "172.18.28.241:2380", "172.18.28.241:2381"},
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	err = cli.Watch(context.Background(), "/test-keep-alive/", func(ctx context.Context, e Event) {
		t.Logf("receive event: key=%v, type=%v", e.Key, e.Type)
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	for i := 0; i < 10000; i++ {
		err = cli.Set(context.Background(), fmt.Sprintf("/test-keep-alive/test_%v", i), "test", 3, true)
		if err != nil {
			t.Fatal(err)
			return
		}
	}

	time.Sleep(time.Second * 60 * 5)
}
