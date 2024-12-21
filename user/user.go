package user

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/types"
	"github.com/junqirao/gocomponents/user/controller"
	"github.com/junqirao/gocomponents/user/logic"
)

type (
	Plugin struct {
		ctx         context.Context
		prefix      string
		middlewares []ghttp.HandlerFunc
		ghttp.Plugin
	}
	PluginOption func(p *Plugin)
)

var (
	// PluginOptionWithContext sets the context of the plugin
	PluginOptionWithContext = func(ctx context.Context) PluginOption {
		return func(p *Plugin) {
			p.ctx = ctx
		}
	}

	// PluginOptionWithPrefix sets the url path prefix of the plugin
	PluginOptionWithPrefix = func(prefix string) PluginOption {
		return func(p *Plugin) {
			p.prefix = prefix
		}
	}

	// PluginOptionWithMiddlewares overwrites the middlewares
	PluginOptionWithMiddlewares = func(middlewares ...ghttp.HandlerFunc) PluginOption {
		return func(p *Plugin) {
			p.middlewares = middlewares
		}
	}
)

func NewPlugin(opt ...PluginOption) ghttp.Plugin {
	p := &Plugin{
		Plugin: types.CommonGoFramePlugin("c_user", "go components user module"),
	}
	for _, accept := range opt {
		accept(p)
	}
	if p.prefix == "" {
		p.prefix = "/user"
	}
	if p.ctx == nil {
		p.ctx = context.Background()
	}
	return p
}

func (p *Plugin) Install(s *ghttp.Server) error {
	s.Group(p.prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(p.middlewares...)

		group.POST("/create", controller.CreateUser)
		group.GET("/exist-username", controller.CheckUsernameExists)
	})

	return logic.CreateAdminIfNotExists(p.ctx)
}
