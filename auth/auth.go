package auth

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/errors/gerror"
)

// GetAppInfo returns secret hashed AppInfo
type GetAppInfo func(ctx context.Context, appId string) (app *AppInfo, err error)

var (
	getFromLocalConfig GetAppInfo = func(ctx context.Context, appId string) (app *AppInfo, err error) {
		value, ok := m.Load(appId)
		if !ok {
			err = errors.New("invalid app id " + appId)
			return
		}
		app = value.(*AppInfo)
		return
	}
)

func auth(ctx context.Context, token string, getFunc GetAppInfo) (app *AppInfo, err error) {
	app = new(AppInfo)
	ts, secret, err := app.FromToken(token)
	if err != nil {
		err = gerror.Wrap(err, "parse token")
		return
	}
	appInfo, err := getFunc(ctx, app.AppId)
	if err != nil {
		err = gerror.Wrap(err, "fetch app info")
		return
	}
	err = appInfo.Check(appInfo.AppKey, ts, secret)
	return
}
