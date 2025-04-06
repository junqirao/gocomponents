package trace

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	uuid "github.com/satori/go.uuid"
)

const (
	HeaderKeyTraceId = "X-Trace-ID"
	ctxKeyTraceId    = "__traceId"
)

// Middleware
// create trace id if request header not set HeaderKeyTraceId
// and add trace id to context value, use GetTraceId to get.
func Middleware(r *ghttp.Request) {
	traceId := r.Header.Get(HeaderKeyTraceId)
	if traceId == "" {
		traceId = uuid.NewV4().String()
		r.Header.Set(HeaderKeyTraceId, traceId)
	}
	r.SetCtxVar(ctxKeyTraceId, traceId)
	r.Middleware.Next()
}

// GetTraceId from context
func GetTraceId(ctx context.Context) string {
	return gconv.String(ctx.Value(ctxKeyTraceId))
}

// CopyTraceInfo from context
func CopyTraceInfo(src context.Context, dst ...context.Context) (ctx context.Context) {
	var newCtx context.Context
	if len(dst) > 0 && dst[0] != nil {
		newCtx = dst[0]
	} else {
		newCtx = context.Background()
	}
	ctx = context.WithValue(newCtx, ctxKeyTraceId, GetTraceId(src))
	return
}
