package audit

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/grace"
)

var (
	Logger  *logger
	modules sync.Map
)

func Init(ctx context.Context) {
	modules = sync.Map{}
	Logger = newLogger(ctx)
	grace.Register(ctx, "FlushAuditLog", func() {
		Sys(ctx, EventShutdown, g.Map{})
		// wait 1s, make sure all logs flushed
		g.Log().Infof(ctx, "waiting 1 seconds for shutdown...")
		time.Sleep(1 * time.Second)
		_ = Logger.Close()
	}, 1000)
	Sys(ctx, EventStartUp, g.Map{})
}

func SetAdaptor(adaptor RecordAdaptor) {
	Logger.SetAdaptor(adaptor)
}

func RegisterModules(ms ...string) {
	for _, module := range ms {
		modules.Store(module, struct{}{})
	}
}

func DeRegisterModules(ms ...string) {
	for _, module := range ms {
		modules.Delete(module)
	}
}

func SupportedModules() []string {
	ms := make([]string, 0)
	modules.Range(func(key, value interface{}) bool {
		ms = append(ms, key.(string))
		return true
	})
	sort.Strings(ms)
	return ms
}
