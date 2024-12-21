package types

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	commonPlugin struct {
		name        string
		description string
	}
)

func (c commonPlugin) Name() string {
	return c.name
}

func (c commonPlugin) Author() string {
	return "go component"
}

func (c commonPlugin) Version() string {
	return "1.0.0"
}

func (c commonPlugin) Description() string {
	return c.description
}

func (c commonPlugin) Install(s *ghttp.Server) error {
	return nil
}

func (c commonPlugin) Remove() error {
	return nil
}

func CommonGoFramePlugin(name string, description string) ghttp.Plugin {
	return &commonPlugin{
		name:        name,
		description: description,
	}
}
