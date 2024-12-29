package gfutil

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	// PluginWrapper plugin wrapper
	PluginWrapper struct {
		prefix         string
		name           string
		author         string
		description    string
		version        string
		middleware     []ghttp.HandlerFunc
		removeHandler  func() error
		installHandler PluginRouterGroupHandler
	}
	// Option of wrapper
	Option                   func(wrapper *PluginWrapper)
	PluginRouterGroupHandler func(group *ghttp.RouterGroup)
)

func NewPlugin(opt ...Option) ghttp.Plugin {
	wrapper := PluginWrapper{}
	for _, o := range opt {
		o(&wrapper)
	}
	return wrapper
}

func (p PluginWrapper) Name() string {
	return p.name
}

func (p PluginWrapper) Author() string {
	return p.author
}

func (p PluginWrapper) Version() string {
	return p.version
}

func (p PluginWrapper) Description() string {
	return p.description
}

func (p PluginWrapper) Install(s *ghttp.Server) error {
	s.Group(p.prefix, func(group *ghttp.RouterGroup) {
		if len(p.middleware) > 0 {
			group.Middleware(p.middleware...)
		}
		if p.installHandler != nil {
			p.installHandler(group)
		}
	})
	return nil
}

func (p PluginWrapper) Remove() error {
	if p.removeHandler != nil {
		return p.removeHandler()
	}
	return nil
}
