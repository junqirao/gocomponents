package meta

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/junqirao/gocomponents/jwt"
	"github.com/junqirao/gocomponents/trace"
)

const (
	ctxKeyMeta = "_meta"
)

type (
	Meta struct {
		User    *User    `json:"user,omitempty"`
		Request *Request `json:"request,omitempty"`
		Server  *Server  `json:"server,omitempty"`
	}

	Request struct {
		Method        string `json:"method"`
		Url           string `json:"url"`
		RemoteAddr    string `json:"remote_addr"`
		WebServerName string `json:"web_server_name"`
		EnterTime     int64  `json:"enter_time"`
		TraceId       string `json:"trace_id"`
	}

	Server struct {
		ServerName string `json:"server_name"`
		HostName   string `json:"host_name"`
		InstanceId string `json:"instance_id"`
	}

	User struct {
		UserId   interface{} `json:"user_id,omitempty"`
		UserName interface{} `json:"user_name,omitempty"`
		UserFrom interface{} `json:"user_from,omitempty"`
		AppId    interface{} `json:"app_id,omitempty"`
	}
)

func (s Server) clone() *Server {
	return &Server{
		ServerName: s.ServerName,
		HostName:   s.HostName,
		InstanceId: s.InstanceId,
	}
}

func Middleware(r *ghttp.Request) {
	meta := &Meta{
		User: &User{
			UserId:   r.GetHeader(jwt.HeaderKeyUserId),
			UserName: r.GetHeader(jwt.HeaderKeyUserName),
			UserFrom: r.GetHeader(jwt.HeaderKeyUserFrom),
		},
		Request: &Request{
			Method:        r.Method,
			Url:           r.URL.String(),
			RemoteAddr:    r.GetRemoteIp(),
			WebServerName: r.Server.GetName(),
			EnterTime:     r.EnterTime.UnixMilli(),
			TraceId:       trace.GetTraceId(r.Context()),
		},
		Server: server.clone(),
	}
	r.SetCtxVar(ctxKeyMeta, meta)
	r.Middleware.Next()
}

func FromCtx(ctx context.Context) *Meta {
	v := ctx.Value(ctxKeyMeta)
	meta := new(Meta)
	_ = gconv.Struct(v, &meta)
	return meta
}

func ServerInfo() Server {
	return *server.clone()
}

func ServerName() string {
	return server.ServerName
}

func HostName() string {
	return server.HostName
}

func InstanceId() string {
	return server.InstanceId
}

func IPV4() string {
	return ipv4
}

func StartedAt() time.Time {
	return startedAt
}
