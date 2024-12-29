package audit

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type logger struct {
	adaptor RecordAdaptor
	ch      chan *Record

	parent *logger
	ctx    context.Context
	module string
}

func newLogger(ctx context.Context, adaptor ...RecordAdaptor) *logger {
	var (
		ad RecordAdaptor = &emptyAdaptor{}
		l  *logger
	)
	if len(adaptor) > 0 && adaptor[0] != nil {
		ad = adaptor[0]
	}
	l = &logger{
		adaptor: ad,
		ch:      make(chan *Record, 100),
		ctx:     ctx,
	}
	go l.cycleHandler()
	return l
}

func (l *logger) SetAdaptor(adapter RecordAdaptor) {
	l.adaptor = adapter
}

func (l *logger) clone() *logger {
	if l.parent != nil {
		return l
	}

	return &logger{
		parent:  l,
		adaptor: l.adaptor,
		ch:      l.ch,
	}
}

func (l *logger) Ctx(ctx context.Context) *logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

func (l *logger) Module(module string) *logger {
	ll := l.clone()
	ll.module = module
	return ll
}

func (l *logger) Log(event string, content ...interface{}) {
	ctx := l.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	l.ch <- &Record{
		Module:    l.module,
		Event:     event,
		Content:   Content(ctx, content...),
		CreatedAt: time.Now(),
	}
	return
}

func (l *logger) cycleHandler() {
	g.Log().Info(l.ctx, "audit start cycle handler")
	defer func() {
		g.Log().Info(l.ctx, "audit cycle handler stopped")
	}()

	for {
		select {
		case <-l.ctx.Done():
			return
		case record, ok := <-l.ch:
			if !ok {
				return
			}
			l.save(record)
		}
	}
}

func (l *logger) Close() error {
	l.ctx.Done()
	close(l.ch)
	for record := range l.ch {
		l.save(record)
	}
	g.Log().Info(l.ctx, "audit cycle handler closed")
	return nil
}

func (l *logger) save(record *Record) {
	if record == nil {
		return
	}
	if record.Content != nil {
		if content, ok := record.Content.(*RecordContent); ok && content != nil {
			if content.Meta != nil && content.Meta.Server != nil {
				record.From = content.Meta.Server.ServiceName
			}
		}
	}
	if err := l.adaptor.Store(l.ctx, record); err != nil {
		g.Log().Errorf(l.ctx, "audit cycle handler error: %v", err)
	}
}
