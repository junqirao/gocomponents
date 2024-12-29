package launcher

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"github.com/junqirao/gocomponents/grace"
	"github.com/junqirao/gocomponents/meta"
	"github.com/junqirao/gocomponents/registry"
)

type (
	// HookFN hook function
	HookFN func(ctx context.Context) error
	// HookTask hook task
	HookTask struct {
		name string
		fn   HookFN
	}
	config struct {
		ctx             context.Context
		init            []*HookTask
		before          []*HookTask
		disableRegistry bool
		updater         func(ctx context.Context) error
	}
)

func NewHookTask(name string, fn func(ctx context.Context) error) *HookTask {
	return &HookTask{name: name, fn: fn}
}

func Launch(blocked func(ctx context.Context), opt ...Option) {
	cfg := new(config)
	for _, o := range opt {
		o(cfg)
	}
	if cfg.ctx == nil {
		cfg.ctx = gctx.GetInitCtx()
	}
	launch(blocked, cfg)
}

func launch(do func(ctx context.Context), cfg *config) {
	ctx := cfg.ctx

	// exec init hooks
	execHooks(ctx, "init", cfg.init)

	// meta init trigger
	g.Log().Infof(ctx, "launcher start service: %s", meta.ServiceName())

	// updater
	if cfg.updater != nil {
		if err := cfg.updater(ctx); err != nil {
			// must update success
			panic(err)
		}
	}

	if !cfg.disableRegistry {
		if err := registry.Init(ctx); err != nil {
			panic(err)
		}
	}

	// exec before hooks
	execHooks(ctx, "before", cfg.before)
	// launch
	if do != nil {
		go do(ctx)
	}

	// exec grace
	grace.GracefulExit(ctx)
}

func execHooks(ctx context.Context, stage string, tasks []*HookTask) {
	g.Log().Infof(ctx, "start execute hook tasks: stage=%s, count=%v", stage, len(tasks))
	execStart := time.Now()
	for _, task := range tasks {
		if task == nil || task.fn == nil {
			continue
		}
		start := time.Now()
		err := task.fn(ctx)
		g.Log().Infof(ctx, "[%s] exec hook task done: cost=%dms, error=%v",
			task.name, time.Now().Sub(start).Milliseconds(), err)
	}
	g.Log().Infof(ctx, "exec hook tasks done: stage=%s, cost=%dms, count=%v",
		stage, time.Now().Sub(execStart).Milliseconds(), len(tasks))
}
