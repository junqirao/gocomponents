package simple_registry

import (
	"fmt"
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
		Prefix            string `json:"prefix"`
		HeartBeatInterval int64  `json:"heart_beat_interval"` // default 3s
	}
)

func (c *Config) check() {
	if c.Prefix == "" {
		c.Prefix = defaultRegistryPrefix
	}
	if c.HeartBeatInterval == 0 {
		c.HeartBeatInterval = defaultHeartBeatInterval
	}
}

func (c *Config) getRegistryPrefix() string {
	return fmt.Sprintf("%sregistry/", c.Prefix)
}
