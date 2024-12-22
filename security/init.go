package security

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	defaultKeypairBits = 2048
	StorageTypeLocal   = "local"
	StorageTypeMysql   = "mysql"
)

type (
	Storage interface {
		StorePublicKey(ctx context.Context, data []byte) error
		StorePrivateKey(ctx context.Context, data []byte) error
		LoadPublicKey(ctx context.Context) error
		LoadPrivateKey(ctx context.Context) error
	}
)

// Init security module from local file or other storage like database
//
// try get keypair from Storage it will auto create if not exist both.
// *notice* make sure every node has the same keypair or maybe cannot
// decode the encrypted data properly cause of keypair mismatch.
func Init(ctx context.Context, s ...Storage) (err error) {
	var storage Storage
	if len(s) == 0 {
		storage = getStorage(ctx)
	} else {
		storage = s[0]
	}

	if err = storage.LoadPrivateKey(ctx); err != nil {
		return
	}
	if err = storage.LoadPublicKey(ctx); err != nil {
		return
	}

	if privateKey == nil && publicKey == nil {
		// auto create
		bits := defaultKeypairBits
		if v, err := g.Cfg().Get(ctx, "security.keypair_bits", defaultKeypairBits); err == nil && v.Int() > 0 {
			bits = v.Int()
		}

		privateKey, publicKey, err = generateKeypair(bits)
		if err == nil {
			g.Log().Infof(ctx, "security module keypair generated bits: %d", bits)
			err = marshalAndStoreCurrentKeypair(ctx, storage)
		}
		return
	} else if privateKey != nil && publicKey != nil {
		return
	}

	return errors.New("incomplete keypair")
}

func getStorage(ctx context.Context) Storage {
	typ := StorageTypeLocal
	v, err := g.Cfg().Get(ctx, "security.storage.type")
	if err == nil {
		if s := v.String(); s != "" {
			typ = s
		}
	}

	g.Log().Infof(ctx, "security module storage loaded type: %s", typ)

	switch typ {
	case StorageTypeMysql:
		return NewMysqlStorage(ctx)
	default:
		return NewLocalStorage()
	}
}

func marshalAndStoreCurrentKeypair(ctx context.Context, storage Storage) (err error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
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
	if err = storage.StorePrivateKey(ctx, bb.Bytes()); err != nil {
		return
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
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
	err = storage.StorePublicKey(ctx, bb.Bytes())
	return
}

func decodePrivateKeyPem(content []byte) (key *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return
	}

	v, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		key = v.(*rsa.PrivateKey)
	}
	return
}

func decodePublicKeyPem(content []byte) (key *rsa.PublicKey, err error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return
	}

	v, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		key = v.(*rsa.PublicKey)
	}
	return
}

func generateKeypair(bits int) (pri *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	if pri, err = rsa.GenerateKey(rand.Reader, bits); err != nil {
		return
	}
	pub = &pri.PublicKey
	return
}
