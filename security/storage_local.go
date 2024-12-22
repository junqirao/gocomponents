package security

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	defaultPublicKeyPath    = "public_key.pem"
	defaultPrivateKeyPath   = "private_key.pem"
	privateKeyConfigPattern = "security.storage.local.private_key"
	publicKeyConfigPattern  = "security.storage.local.public_key"
)

type localStorage struct{}

func NewLocalStorage() Storage {
	return &localStorage{}
}

func (l *localStorage) StorePublicKey(ctx context.Context, data []byte) error {
	return l.store(ctx, defaultPublicKeyPath, publicKeyConfigPattern, data)
}

func (l *localStorage) StorePrivateKey(ctx context.Context, data []byte) error {
	return l.store(ctx, defaultPrivateKeyPath, privateKeyConfigPattern, data)
}

func (l *localStorage) store(ctx context.Context, defPath, pattern string, data []byte) error {
	path := defPath
	v, _ := g.Cfg().Get(ctx, pattern)
	if v.String() != "" {
		path = v.String()
	}
	return gfile.PutBytes(path, data)
}

func (l *localStorage) LoadPublicKey(ctx context.Context) (err error) {
	path := defaultPublicKeyPath
	v, _ := g.Cfg().Get(ctx, publicKeyConfigPattern)
	if v.String() != "" {
		path = v.String()
	}
	content := gfile.GetBytes(path)
	if len(content) > 0 {
		publicKey, err = decodePublicKeyPem(content)
	}
	return
}

func (l *localStorage) LoadPrivateKey(ctx context.Context) (err error) {
	path := defaultPrivateKeyPath
	v, _ := g.Cfg().Get(ctx, privateKeyConfigPattern)
	if v.String() != "" {
		path = v.String()
	}
	content := gfile.GetBytes(path)
	if len(content) > 0 {
		privateKey, err = decodePrivateKeyPem(content)
	}
	return
}
