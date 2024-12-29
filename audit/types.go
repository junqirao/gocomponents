package audit

import (
	"context"
	"io"
	"time"

	"git.lumosc.com/lumosc/core/meta"
)

type ILogger interface {
	io.Closer
	Ctx(ctx context.Context) *logger
	Module(module string) *logger
	Log(event string, content ...interface{})
}

type Record struct {
	Module    string      `json:"module"`     // module
	Event     string      `json:"event"`      // event name
	From      string      `json:"from"`       // identity
	Content   interface{} `json:"content"`    // content
	CreatedAt time.Time   `json:"created_at"` // created_at
}

type RecordContent struct {
	Meta  *meta.Meta  `json:"meta"`
	Value interface{} `json:"value"`
	Error string      `json:"error"`
}

func Content(ctx context.Context, v ...interface{}) *RecordContent {
	var (
		value  interface{}
		errMsg string
	)

	if len(v) > 0 && v[0] != nil {
		value = v[0]
	}
	if len(v) > 1 && v[1] != nil {
		if e, ok := v[1].(error); ok && e != nil {
			errMsg = e.Error()
		}
	}
	return &RecordContent{
		Meta:  meta.FromCtx(ctx),
		Value: value,
		Error: errMsg,
	}
}
