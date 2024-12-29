package grace

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"

	"github.com/gogf/gf/v2/frame/g"
)

type Handler func()

type handlerWithPriority struct {
	Name     string
	Handler  Handler
	Priority int
}

type handlers []*handlerWithPriority

func (h handlers) Len() int {
	return len(h)
}

func (h handlers) Less(i, j int) bool {
	return h[i].Priority < h[j].Priority
}

func (h handlers) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h handlers) append(handler ...*handlerWithPriority) handlers {
	return append(h, handler...)
}

var m = sync.Map{}

func Register(ctx context.Context, name string, handler Handler, priority ...int) {
	g.Log().Infof(ctx, "grace register handler: %s", name)
	p := 0
	if len(priority) > 0 {
		p = priority[0]
	}
	m.Store(name, &handlerWithPriority{
		Name:     name,
		Handler:  handler,
		Priority: p,
	})
}

func UnRegister(ctx context.Context, name string) {
	g.Log().Infof(ctx, "grace unregister handler: %s", name)
	m.Delete(name)
}

func ExecAndExit(ctx context.Context) {
	hp := handlers{}
	m.Range(func(name, handler any) bool {
		hp = hp.append(handler.(*handlerWithPriority))
		return true
	})
	sort.Sort(hp)
	for _, handler := range hp {
		handler.Handler()
		g.Log().Infof(ctx, "exec grace handler %s done.", handler.Name)
	}
	g.Log().Infof(ctx, "grace handler all executed: %d", len(hp))
	g.Log().Info(ctx, "shutdown.")
	os.Exit(0)
}

func GracefulExit(ctx context.Context) {
	osc := make(chan os.Signal)
	signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)

	// block and wait for signal
	s := <-osc
	g.Log().Infof(ctx, "receive stop sig: %s", s)
	ExecAndExit(ctx)
}
