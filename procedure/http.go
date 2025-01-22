package procedure

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type (
	GoFrameHTTPInputHandler struct{}
)

func (g GoFrameHTTPInputHandler) HandleInput(ctx context.Context, proto *Proto) (input map[string]any, err error) {
	input = make(map[string]any)
	req := ghttp.RequestFromCtx(ctx)
	if req == nil {
		err = errors.New("not go frame http request")
		return
	}
	reqMap := req.GetRequestMap()

	for _, parameter := range proto.Parameters {
		val := gvar.New(nil)
		switch parameter.Meta.Get("from").String() {
		case "header":
			val.Set(req.GetHeader(parameter.Name))
		case "query":
			val = req.GetQuery(parameter.Name)
		case "path":
			val = req.GetRouter(parameter.Name)
		default:
			if v, ok := reqMap[parameter.Name]; ok {
				val.Set(v)
			} else {
				// try to get header
				// reqMap not include header
				val.Set(req.GetHeader(parameter.Name))
			}
		}
		if !val.IsNil() && !val.IsEmpty() {
			input[parameter.Name] = val.Interface()
		}
	}
	return
}

func BindGHTTPRouter(ctx context.Context, group *ghttp.RouterGroup, lc LifeCycle, responseHandler func(r *ghttp.Request, res any, err error), proto ...*Proto) (err error) {
	for _, p := range proto {
		p.LifeCycle = lc
		if err = p.Check(ctx); err != nil {
			glog.Warningf(ctx, "proto %s check error: %v", p.Name, err.Error())
			return
		}
		var (
			path     = p.Meta.Get("path").String()
			method   = strings.ToUpper(p.Meta.Get("method").String())
			handler  = wrapGoFrameHttpHandler(p, responseHandler)
			bindFunc func(pattern string, object interface{}, params ...interface{}) *ghttp.RouterGroup
		)

		if path == "" {
			err = fmt.Errorf("proto %s path not set, skip", p.Name)
			return
		}

		switch method {
		case http.MethodGet:
			bindFunc = group.GET
		case http.MethodPost:
			bindFunc = group.POST
		case http.MethodPut:
			bindFunc = group.PUT
		case http.MethodDelete:
			bindFunc = group.DELETE
		case http.MethodPatch:
			bindFunc = group.PATCH
		case http.MethodHead:
			bindFunc = group.HEAD
		case http.MethodOptions:
			bindFunc = group.OPTIONS
		default:
			glog.Warningf(ctx, "proto %s method %s not supported, use ALL", p.Name, method)
			bindFunc = group.ALL
		}
		bindFunc(path, handler)
		glog.Infof(ctx, "bind proto %s path %s method %s done", p.Name, path, method)
	}
	return
}

func wrapGoFrameHttpHandler(p *Proto, responseHandler func(r *ghttp.Request, res any, err error)) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		out, err := p.Execute(r.Context())
		responseHandler(r, out, err)
	}
}
