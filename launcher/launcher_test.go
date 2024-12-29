package launcher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/grace"
)

func TestLaunch(t *testing.T) {
	before := make([]*HookTask, 0)
	for i := 0; i < 5; i++ {
		index := i
		task := NewHookTask(fmt.Sprintf("task-%d", index), func(ctx context.Context) error {
			time.Sleep(time.Millisecond * 500)
			return nil
		})
		before = append(before, task)
	}

	grace.Register(context.Background(), "after1", func() {
		fmt.Println("do after1...")
	})
	grace.Register(context.Background(), "after2", func() {
		fmt.Println("do after2...")
	})

	Launch(func(ctx context.Context) {
		g.Log().Infof(ctx, "do something, 2s")
		time.Sleep(time.Second * 2)
		g.Log().Infof(ctx, "done")
	},
		WithBeforeTasks(before...),
		DisableRegistry(false))
}
