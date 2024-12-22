package security

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
)

func TestLoadFromLocal(t *testing.T) {
	err := Load(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	printKey := func(typeName string, bs []byte) {
		bb := &bytes.Buffer{}
		_ = pem.Encode(bb, &pem.Block{
			Type:  typeName,
			Bytes: bs,
		})

		fmt.Println(bb.String())
	}

	pri, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	printKey("PRIVATE KEY", pri)
	fmt.Println()
	fmt.Println()
	pub, _ := x509.MarshalPKIXPublicKey(publicKey)
	printKey("PUBLIC KEY", pub)
}

func TestEncryptDecrypt(t *testing.T) {
	err := Load(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	data := "hello world"
	encrypted, err := Encrypt(data)
	if err != nil {
		t.Fatal(err)
		return
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
		return
	}

	if decrypted != data {
		t.Fatal("data not match")
	}
}

func TestGetPublicKeyPem(t *testing.T) {
	err := Load(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(GetPublicKeyPem())
}
