package kvdb

import (
	"crypto/tls"
)

type (
	// Config for registry
	Config struct {
		// common
		Endpoints []string `json:"endpoints"`
		Username  string   `json:"username"`
		Password  string   `json:"password"`
		// Etcd tls
		Tls *TlsConfig `json:"tls"`
	}
	// StorageConfig for storage module
	StorageConfig struct {
		Separator string `json:"separator"`
	}
	// TlsConfig ...
	TlsConfig struct {
		InsecureSkipVerify bool `json:"insecure_skip_verify"`
	}
)

func (c Config) tlsConfig() *tls.Config {
	if c.Tls == nil {
		return nil
	}
	return &tls.Config{
		InsecureSkipVerify: c.Tls.InsecureSkipVerify,
	}
}
