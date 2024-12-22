package security

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	defaultPublicKeyPath    = "./"
	defaultPrivateKeyPath   = "./"
	privateKeyConfigPattern = "security.storage.local.private_key"
	publicKeyConfigPattern  = "security.storage.local.public_key"
)

type localStorage struct{}

func NewLocalStorage() Storage {
	return &localStorage{}
}

func (l *localStorage) StorePublicKey(ctx context.Context, name string, data []byte) error {
	return l.store(ctx, defaultPublicKeyPath, publicKeyConfigPattern, l.getFileName(name, "public_key"), data)
}

func (l *localStorage) StorePrivateKey(ctx context.Context, name string, data []byte) error {
	return l.store(ctx, defaultPrivateKeyPath, privateKeyConfigPattern, l.getFileName(name, "s_private_key"), data)
}

func (l *localStorage) store(ctx context.Context, defPath, pattern, name string, data []byte) error {
	path := defPath
	v, _ := g.Cfg().Get(ctx, pattern)
	if v.String() != "" {
		path = v.String()
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	path += name
	return gfile.PutBytes(path, data)
}

func (l *localStorage) LoadPublicKey(ctx context.Context, name string) (publicKey *rsa.PublicKey, err error) {
	content := l.getFileBytes(ctx, l.getFileName(name, "public_key"), defaultPublicKeyPath, publicKeyConfigPattern)
	if len(content) > 0 {
		publicKey, err = decodePublicKeyPem(content)
	}
	return
}

func (l *localStorage) LoadPrivateKey(ctx context.Context, name string) (privateKey *rsa.PrivateKey, err error) {
	content := l.getFileBytes(ctx, l.getFileName(name, "private_key"), defaultPrivateKeyPath, privateKeyConfigPattern)
	if len(content) > 0 {
		privateKey, err = decodePrivateKeyPem(content)
	}
	return
}

func (l *localStorage) getFileBytes(ctx context.Context, name, defPath, cfgPattern string) []byte {
	path := defPath
	v, _ := g.Cfg().Get(ctx, cfgPattern)
	if v.String() != "" {
		path = v.String()
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	path += name
	return gfile.GetBytes(path)
}

func (l *localStorage) getFileName(name string, typ string) string {
	return fmt.Sprintf("%s_%s.pem", name, typ)
}
