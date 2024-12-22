package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/gogf/gf/v2/util/gconv"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// Encrypt using rsa
func Encrypt(data interface{}) (s string, err error) {
	if publicKey == nil {
		return "", errors.New("public key not found")
	}
	bs, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(gconv.String(data)))
	if err != nil {
		return
	}

	s = base64.StdEncoding.EncodeToString(bs)
	return
}

// Decrypt using rsa
func Decrypt(data interface{}) (s string, err error) {
	if privateKey == nil {
		return "", errors.New("private key not found")
	}
	raw, err := base64.StdEncoding.DecodeString(gconv.String(data))
	if err != nil {
		return
	}

	bs, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, raw)
	if err != nil {
		return
	}

	s = string(bs)
	return
}

// GetPublicKeyPem get public key pem
func GetPublicKeyPem() string {
	bs, _ := x509.MarshalPKIXPublicKey(publicKey)
	bb := &bytes.Buffer{}
	_ = pem.Encode(bb, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bs,
	})
	return bb.String()
}
