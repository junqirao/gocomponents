package security

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type (
	Provider struct {
		name       string
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		storage    Storage
	}
)

func NewProvider(ctx context.Context, name, storageType string) (p *Provider, err error) {
	p = &Provider{
		name:    name,
		storage: getStorageByType(ctx, storageType),
	}
	err = p.setup(ctx)
	return
}

func (p *Provider) setup(ctx context.Context) (err error) {
	if p.privateKey, err = p.storage.LoadPrivateKey(ctx, p.name); err != nil {
		return
	}
	if p.publicKey, err = p.storage.LoadPublicKey(ctx, p.name); err != nil {
		return
	}

	if p.privateKey == nil && p.publicKey == nil {
		// auto create
		bits := defaultKeypairBits
		if v, err := g.Cfg().Get(ctx, "security.keypair_bits", defaultKeypairBits); err == nil && v.Int() > 0 {
			bits = v.Int()
		}

		p.privateKey, p.publicKey, err = generateKeypair(bits)
		if err == nil {
			g.Log().Infof(ctx, "security module instance[%s] keypair generated bits: %d", p.name, bits)
			err = p.marshalAndStoreCurrentKeypair(ctx, p.storage)
		}
		return
	} else if p.privateKey != nil && p.publicKey != nil {
		return
	}

	return errors.New("incomplete keypair")
}

// Encrypt using rsa
func (p *Provider) Encrypt(data interface{}) (s string, err error) {
	if p.publicKey == nil {
		return "", errors.New("public key not found")
	}
	bs, err := rsa.EncryptPKCS1v15(rand.Reader, p.publicKey, []byte(gconv.String(data)))
	if err != nil {
		return
	}

	s = base64.StdEncoding.EncodeToString(bs)
	return
}

// Decrypt using rsa
func (p *Provider) Decrypt(data interface{}) (s string, err error) {
	if p.privateKey == nil {
		return "", errors.New("private key not found")
	}
	raw, err := base64.StdEncoding.DecodeString(gconv.String(data))
	if err != nil {
		return
	}

	bs, err := rsa.DecryptPKCS1v15(rand.Reader, p.privateKey, raw)
	if err != nil {
		return
	}

	s = string(bs)
	return
}

// GetPublicKeyPem get public key pem
func (p *Provider) GetPublicKeyPem() string {
	bs, _ := x509.MarshalPKIXPublicKey(p.publicKey)
	bb := &bytes.Buffer{}
	_ = pem.Encode(bb, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bs,
	})
	return bb.String()
}

func (p *Provider) marshalAndStoreCurrentKeypair(ctx context.Context, storage Storage) (err error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(p.privateKey)
	if err != nil {
		return
	}
	bb := &bytes.Buffer{}
	if err = pem.Encode(bb, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}); err != nil {
		return
	}
	if err = storage.StorePrivateKey(ctx, p.name, bb.Bytes()); err != nil {
		return
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(p.publicKey)
	if err != nil {
		return
	}
	bb = &bytes.Buffer{}
	if err = pem.Encode(bb, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}); err != nil {
		return
	}
	err = storage.StorePublicKey(ctx, p.name, bb.Bytes())
	return
}
