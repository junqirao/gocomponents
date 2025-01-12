package kvdb

import (
	"crypto/tls"
)

const (
	defaultIdentitySeparator = "/"
	defaultStoragePrefix     = "/storage/"
)

type (
	// Config for etcd
	Config struct {
		// common
		Endpoints []string `json:"endpoints"`
		Username  string   `json:"username"`
		Password  string   `json:"password"`
		// tls
		Tls *TlsConfig `json:"tls"`
		// separator for key in database, default "/"
		Separator string `json:"separator"`
		// storage module config
		Storage StorageConfig `json:"storage"`
	}
	StorageConfig struct {
		// key prefix, default "/storage/"
		// start with "/",and end with "/" in etcd
		Prefix string `json:"prefix"`
	}

	// TlsConfig ...
	TlsConfig struct {
		InsecureSkipVerify bool `json:"insecure_skip_verify"`
	}
)

func (c *Config) tlsConfig() *tls.Config {
	if c.Tls == nil {
		return nil
	}
	return &tls.Config{
		InsecureSkipVerify: c.Tls.InsecureSkipVerify,
	}
}

func (c *Config) check() {
	if c.Separator == "" {
		c.Separator = defaultIdentitySeparator
	}
	if c.Storage.Prefix == "" {
		c.Storage.Prefix = defaultStoragePrefix
	}
}
