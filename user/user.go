package user

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/types"
	"github.com/junqirao/gocomponents/user/controller"
	"github.com/junqirao/gocomponents/user/logic"
)

const (
	tableName      = "c_user"
	createTableDDL = `
create table if not exists c_user
(
    id            varchar(50)            not null primary key,
    username      varchar(20)            not null,
    password      varchar(200)           not null,
    created_at    datetime               null,
    updated_at    datetime               null,
    administrator tinyint(1)  default 0  null,
    source        varchar(20) default '' null,
    status        tinyint     default 0  null,
    extra         json                   null,
    unique index c_user_uk_username (username)
);
`
)

type (
	Plugin struct {
		db          string
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

	// PluginOptionWithDBGroupName sets the database group name
	PluginOptionWithDBGroupName = func(db string) PluginOption {
		return func(p *Plugin) {
			p.db = db
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
		p.prefix = "/"
	}
	if p.ctx == nil {
		p.ctx = context.Background()
	}
	return p
}

func (p *Plugin) Install(s *ghttp.Server) (err error) {
	if err = p.createTableIfNotExists(p.ctx); err != nil {
		return
	}
	s.Group(p.prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(p.middlewares...)

		group.Group("/user", func(group *ghttp.RouterGroup) {
			// no need login
			group.POST("/register", controller.CreateUser)
			group.POST("/exist-username", controller.CheckUsernameExists)
			group.POST("/login", controller.Login)

		})

		group.Group("/security", func(group *ghttp.RouterGroup) {
			group.GET("/transport/public-key", controller.GetPublicKeyPem)
		})

	})

	return logic.CreateAdminIfNotExists(p.ctx)
}

func (p *Plugin) createTableIfNotExists(ctx context.Context) (err error) {
	tables, err := g.DB().Ctx(ctx).Tables(ctx)
	if err != nil {
		return
	}
	for _, table := range tables {
		if table == tableName {
			return
		}
	}
	_, err = g.DB(p.db).Ctx(ctx).Exec(ctx, createTableDDL)
	return
}
