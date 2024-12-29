package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	defaultRegistryPrefix    = "/default-registry-service/"
	defaultHeartBeatInterval = 3
	defaultIdentitySeparator = "/"
	defaultPort              = 8000
)

type (
	// Config for registry
	Config struct {
		Instance          *Instance `json:"instance"`
		Prefix            string    `json:"prefix"`
		HeartBeatInterval int64     `json:"heart_beat_interval"` // default 3s
	}
)

func (c *Config) check() {
	if c.Prefix == "" {
		c.Prefix = defaultRegistryPrefix
	}
	if c.HeartBeatInterval == 0 {
		c.HeartBeatInterval = defaultHeartBeatInterval
	}
	if c.Instance.Port == 0 {
		// try to get server.address
		v, err := g.Cfg().Get(context.Background(), "server.address")
		if err == nil {
			parts := strings.Split(v.String(), ":")
			if len(parts) == 0 {
				return
			}
			c.Instance.Port = gconv.Int(parts[len(parts)-1])
		}
	}
}

func (c *Config) getRegistryPrefix() string {
	return fmt.Sprintf("%sregistry/", c.Prefix)
}
