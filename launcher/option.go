package launcher

import (
	"context"
	"fmt"

	"github.com/junqirao/gocomponents/kvdb"
	"github.com/junqirao/gocomponents/meta"
	"github.com/junqirao/gocomponents/updater"
)

// Option ...
type Option func(cfg *config)

// NewOptions of launcher
func NewOptions(op ...Option) []Option {
	return append(make([]Option, 0), op...)
}

// WithContext set context
func WithContext(ctx context.Context) Option {
	return func(cfg *config) {
		cfg.ctx = ctx
	}
}

// DisableRegistry disable registry
func DisableRegistry(disable bool) Option {
	return func(cfg *config) {
		cfg.disableRegistry = disable
	}
}

// WithBeforeTasks add before tasks
func WithBeforeTasks(tasks ...*HookTask) Option {
	return func(cfg *config) {
		m := map[string]int{}
		for i, task := range cfg.before {
			m[task.name] = i
		}
		for _, task := range tasks {
			if index, ok := m[task.name]; ok {
				cfg.before[index] = task
			} else {
				cfg.before = append(cfg.before, task)
			}
		}
	}
}

// WithInitTasks add init tasks
func WithInitTasks(tasks ...*HookTask) Option {
	return func(cfg *config) {
		m := map[string]int{}
		for i, task := range cfg.init {
			m[task.name] = i
		}
		for _, task := range tasks {
			if index, ok := m[task.name]; ok {
				cfg.init[index] = task
			} else {
				cfg.init = append(cfg.init, task)
			}
		}
	}
}

func WithConcurrencyUpdater(
	adaptor updater.RecordAdaptor,
	functions ...*updater.FuncInfo) Option {
	return func(cfg *config) {
		cfg.updater = func(ctx context.Context) (err error) {
			mu, err := kvdb.NewMutex(ctx, fmt.Sprintf("updater_exec_%s", meta.ServiceName()))
			if err != nil {
				return
			}
			return updater.ConcurrencyUpdate2Latest(ctx, adaptor, mu, functions...)
		}
	}
}

func WithUpdater(
	adaptor updater.RecordAdaptor,
	functions ...*updater.FuncInfo) Option {
	return func(cfg *config) {
		cfg.updater = func(ctx context.Context) error {
			return updater.Update2Latest(ctx, adaptor, functions...)
		}
	}
}
