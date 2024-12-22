package security

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"sync"
)

const (
	defaultKeypairBits = 2048
	StorageTypeLocal   = "local"
	StorageTypeMysql   = "mysql"
)

var (
	instances = sync.Map{} // name : *Provider
)

type (
	Storage interface {
		StorePublicKey(ctx context.Context, name string, data []byte) error
		StorePrivateKey(ctx context.Context, name string, data []byte) error
		LoadPublicKey(ctx context.Context, name string) (publicKey *rsa.PublicKey, err error)
		LoadPrivateKey(ctx context.Context, name string) (privateKey *rsa.PrivateKey, err error)
	}
)

func getStorageByType(ctx context.Context, typ string) Storage {
	switch typ {
	case StorageTypeMysql:
		return NewMysqlStorage(ctx)
	default:
		return NewLocalStorage()
	}
}

// GetProvider of security module, if not exists create with specified
// type (using storage type "mysql" as default).
//
// try get keypair from Storage it will auto create if not exist both.
// *notice* make sure every node has the same keypair or maybe cannot
// decode the encrypted data properly cause of keypair mismatch.
func GetProvider(ctx context.Context, name string, typ ...string) (p *Provider, err error) {
	value, ok := instances.Load(name)
	if ok {
		p = value.(*Provider)
		return
	}

	st := StorageTypeMysql
	if len(typ) > 0 {
		st = typ[0]
	}
	if p, err = NewProvider(ctx, name, st); err != nil {
		return
	}

	instances.Store(name, p)
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
