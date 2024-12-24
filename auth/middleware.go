package auth

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/response"
)

var Middleware = MiddlewareWithCheckFunc(getFromLocalConfig)

func MiddlewareWithCheckFunc(getFunc GetAppInfo) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		header := r.GetHeader("Authorization")
		if header == "" {
			response.WriteJSON(r, response.CodeUnauthorized.WithDetail("basic token required"))
			return
		}

		app, err := auth(r.Context(), header, getFunc)
		if err != nil {
			response.WriteJSON(r, response.CodeUnauthorized.WithMessage(err.Error()))
			return
		}

		r.Header.Set(HeaderKeyAppId, app.AppId)
		r.Header.Set(HeaderKeyAppKey, app.AppKey)
		r.Middleware.Next()
	}
}
