package updater

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/junqirao/gocomponents/meta"
)

const (
	FuncTypeRaw = 0
	FuncTypeSql = 1
)

const (
	ExecStatusFailed  = 0
	ExecStatusSuccess = 1
)

var (
	ErrExecuteTimeout = errors.New("execute timeout")
)

type (
	// FN ...
	FN = func(ctx context.Context) (err error)
	// FuncInfo ...
	FuncInfo struct {
		FuncConfig
		Name string
		Type int
		FN   FN
	}
	// FuncConfig function config
	FuncConfig struct {
		// must execute success,
		// if true and execute failed, will break and return error,
		// default false
		must bool
		// can retry in next updater procedure
		// default false
		retry bool
		// function execute timeout
		// default 0, unlimited
		timeout time.Duration
	}
)

func NewFuncConfig() *FuncConfig {
	return &FuncConfig{}
}

func (c *FuncConfig) Must(b ...bool) *FuncConfig {
	if len(b) > 0 {
		c.must = b[0]
	} else {
		c.must = true
	}
	return c
}

func (c *FuncConfig) Retry(b ...bool) *FuncConfig {
	if len(b) > 0 {
		c.retry = b[0]
	} else {
		c.retry = true
	}
	return c
}

func (c *FuncConfig) Timeout(d time.Duration) *FuncConfig {
	c.timeout = d
	return c
}

// NewFunc constructor
func NewFunc(name string, fn FN, cfg *FuncConfig, typ ...int) *FuncInfo {
	t := FuncTypeRaw
	if len(typ) > 0 {
		t = typ[0]
	}
	return &FuncInfo{
		Name:       fmt.Sprintf("%s/%s", meta.ServerName(), name),
		FN:         fn,
		Type:       t,
		FuncConfig: *cfg,
	}
}

func (i FuncInfo) Exec(ctx context.Context) (cost int64, err error) {
	if i.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, i.timeout)
		defer cancel()
	}

	var (
		done = make(chan struct{})
		now  = time.Now()
	)
	defer close(done)

	f := func() {
		err = i.FN(ctx)
		done <- struct{}{}
	}

	go f()

	select {
	case <-ctx.Done():
		err = ErrExecuteTimeout
	case <-done:
	}

	cost = time.Since(now).Milliseconds()
	g.Log().Infof(ctx, "updater func %s executed cost %dms, error=%v", i.Name, cost, err)
	return
}

func Update2Latest(ctx context.Context, adaptor RecordAdaptor, functions ...*FuncInfo) (err error) {
	records, err := adaptor.Load(ctx, &RecordQueryParams{})
	if err != nil {
		g.Log().Errorf(ctx, "updater load record error: %v", err)
		return
	}

	// already executed
	rm := make(map[string]*Record) // record_name : *Record
	for _, record := range records.Records {
		rm[record.Name] = record
	}

	var (
		start   = time.Now()
		success = 0
		failed  = 0
	)
	defer func() {
		g.Log().Infof(ctx, "updater execute done: cost=%dms, success=%d, failed=%d, error=%v",
			time.Now().Sub(start).Milliseconds(), success, failed, err)
	}()

	for _, fi := range functions {
		r, ok := rm[fi.Name]
		if ok {
			if r.Status == ExecStatusSuccess ||
				(r.Status == ExecStatusFailed && !fi.retry) {
				// skip already executed, and cannot retry
				g.Log().Infof(ctx, "updater record %s already executed, skip.", fi.Name)
				continue
			}
		}
		record := &Record{
			Name:      fi.Name,
			Type:      fi.Type,
			CreatedAt: gconv.String(time.Now().UnixMilli()),
		}

		record.Cost, err = fi.Exec(ctx)
		if err != nil {
			g.Log().Warningf(ctx, "exec function %s error: %v", fi.Name, err)
			if !fi.must {
				err = nil
			}
			failed++
		} else {
			record.Status = ExecStatusSuccess
			success++
		}

		if err := adaptor.Store(ctx, record); err != nil {
			g.Log().Errorf(ctx, "updater store record error: %v", err)
			return err
		}

		if fi.must && err != nil {
			return
		}
	}
	return
}

// ConcurrencyUpdate2Latest use distributed lock to update latest.
// see kvdb.NewMutex, if you use sync.Locker locally, it may cause
// consistency problems in distributed environment.
func ConcurrencyUpdate2Latest(ctx context.Context, adaptor RecordAdaptor, mu sync.Locker, functions ...*FuncInfo) (err error) {
	mu.Lock()
	defer mu.Unlock()
	return Update2Latest(ctx, adaptor, functions...)
}
