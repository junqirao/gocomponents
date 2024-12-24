package auth

import (
	"testing"
	"time"
)

func TestBasicToken(t *testing.T) {
	info := &AppInfo{
		AppId:     "test",
		AppKey:    "12345678900",
		AppSecret: "test123456",
	}
	m.Store(info.AppId, info)

	token := info.Token(time.Now().UnixMilli())
	t.Logf("generate token:\n %s", token)
	a := new(AppInfo)
	ts, secret, err := a.FromToken(token)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err = info.Check(a.AppKey, ts, secret); err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("ts=%v, secret=%v", ts, secret)
	t.Logf("app info: %+v", a)
}
