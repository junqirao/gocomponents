package grace

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogf/gf/v2/frame/g"
)

func Init(ctx context.Context) {
	go func() {
		osc := make(chan os.Signal, 1)
		signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)
		s := <-osc
		g.Log().Infof(ctx, "receive stop sig: %s", s)
	}()
}
