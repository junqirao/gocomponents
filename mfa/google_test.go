package mfa

import (
	"testing"
)

func TestGoogleMFA(t *testing.T) {
	authenticator := NewGoogleAuthenticator(6, 16)
	secret, err := authenticator.CreateSecret()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("secret: %s", secret)
	code, err := authenticator.GetCode(secret)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("code: %s", code)
	if !authenticator.VerifyCode(secret, code, 1, 0) {
		t.Fatal("verify failed")
		return
	}

	t.Log("verify success")
	qrcode, err := authenticator.GenerateQRCode("test", secret)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("qrcode: %s", qrcode)
}
