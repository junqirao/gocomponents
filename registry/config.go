package registry

import (
	"fmt"
)

const (
	defaultName              = "default"
	defaultHeartBeatInterval = 10
	defaultIdentitySeparator = "/"
	defaultPort              = 8000
)

type (
	// Config for registry
	Config struct {
		Instance          *Instance `json:"instance"`
		Name              string    `json:"name"`
		HeartBeatInterval int64     `json:"heart_beat_interval"` // default 10s
		MaximumRetry      int       `json:"maximum_retry"`       // default 20
	}
)

func (c *Config) check() {
	if c.Name == "" {
		c.Name = defaultName
	}
	if c.HeartBeatInterval == 0 {
		c.HeartBeatInterval = defaultHeartBeatInterval
	}
	if c.MaximumRetry == 0 {
		c.MaximumRetry = 20
	}
}

func (c *Config) getRegistryPrefix() string {
	return fmt.Sprintf("/registry/%s/", c.Name)
}
