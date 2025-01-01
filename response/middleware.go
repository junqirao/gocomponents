package response

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	DataHandler = func(ctx context.Context, res any) any
)

func Middleware(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err  = r.GetError()
		ec   = CodeFromError(err)
		data = r.GetHandlerResponse()
	)
	if ec == nil {
		ec = DefaultSuccess()
	}

	WriteData(r, ec, data)
}

func MiddlewareWithDataHandler(hs ...DataHandler) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		r.Middleware.Next()

		if r.Response.BufferLength() > 0 {
			return
		}

		var (
			err  = r.GetError()
			ec   = CodeFromError(err)
			data = r.GetHandlerResponse()
		)

		for _, handler := range hs {
			if handler == nil {
				continue
			}
			data = handler(r.Context(), data)
		}

		if ec == nil {
			ec = DefaultSuccess()
		}

		WriteData(r, ec, data)
	}
}
