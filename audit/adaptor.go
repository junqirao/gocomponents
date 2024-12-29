package audit

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"git.lumosc.com/lumosc/core/response"
)

type RecordQueryParams struct {
	Module   string     `json:"module"`
	Event    string     `json:"event"`
	From     string     `json:"from"`
	Start    *time.Time `json:"start"`
	End      *time.Time `json:"end"`
	PageSize int        `json:"page_size"`
	PageNum  int        `json:"page_num"`
}

type RecordQueryResult struct {
	Records []*Record `json:"records"`
	Total   int       `json:"total"`
}

type RecordAdaptor interface {
	Store(ctx context.Context, record *Record) (err error)
	Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error)
}

type emptyAdaptor struct{}

func (a emptyAdaptor) Store(ctx context.Context, record *Record) (err error) {
	g.Log().Cat("audit").Info(ctx, gconv.String(record))
	return
}
func (a emptyAdaptor) Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error) {
	err = response.CodeNotFound.WithDetail("using empty adaptor")
	return
}
