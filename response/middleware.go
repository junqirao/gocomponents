package response

import (
	"github.com/gogf/gf/v2/net/ghttp"
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
