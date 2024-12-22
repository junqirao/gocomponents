package security

import (
	"context"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	var ps []*Provider

	p1, err := GetProvider(context.Background(), "test_mysql", StorageTypeMysql)
	if err != nil {
		t.Fatal(err)
		return
	}
	ps = append(ps, p1)
	p2, err := GetProvider(context.Background(), "test_local", StorageTypeLocal)
	if err != nil {
		t.Fatal(err)
		return
	}
	ps = append(ps, p2)

	data := "hello world"
	for _, p := range ps {
		t.Log(p.GetPublicKeyPem())
		encrypted, err := p.Encrypt(data)
		if err != nil {
			t.Fatal(err)
			return
		}

		decrypted, err := p.Decrypt(encrypted)
		if err != nil {
			t.Fatal(err)
			return
		}

		if decrypted != data {
			t.Fatal("data not match")
		}
	}
}
