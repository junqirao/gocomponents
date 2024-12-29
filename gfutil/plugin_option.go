package gfutil

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func WithPrefix(prefix string) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.prefix = prefix
	}
}

func WithName(name string) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.name = name
	}
}

func WithAuthor(author string) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.author = author
	}
}

func WithDescription(description string) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.description = description
	}
}

func WithVersion(version string) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.version = version
	}
}

func WithMiddleware(middleware ...ghttp.HandlerFunc) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.middleware = middleware
	}
}

func WithRemoveHandler(removeHandler func() error) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.removeHandler = removeHandler
	}
}

func WithInstallHandler(i PluginRouterGroupHandler) Option {
	return func(wrapper *PluginWrapper) {
		wrapper.installHandler = i
	}
}
