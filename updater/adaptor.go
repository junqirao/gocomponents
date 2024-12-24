package updater

import (
	"context"
)

type (
	Record struct {
		Name      string `json:"name"`       // unique name+type
		Type      int    `json:"type"`       // 0: function, 1: sql
		Status    int    `json:"status"`     // 0: failed, 1: success
		Cost      int64  `json:"cost"`       // cost time in ms
		CreatedAt string `json:"created_at"` // created_at
	}
	RecordQueryParams struct {
		Name *string
		Type *int
	}
	RecordQueryResult struct {
		Records []*Record `json:"records"`
		Total   int       `json:"total"`
	}
	RecordAdaptor interface {
		Store(ctx context.Context, record *Record) (err error)
		Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error)
	}
)
