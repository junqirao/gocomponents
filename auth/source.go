package auth

import (
	"context"
)

type (
	DataSource interface {
		List(ctx context.Context, params *ListParams) (res *ListResult, err error)
		Store(ctx context.Context, app *AppFullInfo) (err error)
		Delete(ctx context.Context, appId string) (err error)
		FindOne(ctx context.Context, appId string) (app *AppFullInfo, err error)
	}
	ListParams struct {
		Name     string `json:"name"`
		AppId    string `json:"app_id"`
		Fuzzy    bool   `json:"fuzzy"`
		PageSize int    `json:"page_size"`
		PageNum  int    `json:"page_num"`
	}
	ListResult struct {
		List  []*AppFullInfo `json:"list"`
		Total int            `json:"total"`
	}
)
