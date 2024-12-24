package dkv

import (
	"crypto/tls"
)

type Config struct {
	Endpoints []string   `json:"endpoints"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Tls       *TlsConfig `json:"tls"`
}

func (c Config) tlsConfig() *tls.Config {
	if c.Tls == nil {
		return nil
	}
	return &tls.Config{
		InsecureSkipVerify: c.Tls.InsecureSkipVerify,
	}
}

type TlsConfig struct {
	InsecureSkipVerify bool
}
