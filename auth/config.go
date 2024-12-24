package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	cfg *Config
	m   = sync.Map{}
)

func SetConfig(config *Config) {
	cfg = config
	for _, app := range cfg.Apps {
		m.Store(app.AppId, app)
	}
}

type Config struct {
	Apps []*AppInfo `json:"apps"`
}

type AppInfo struct {
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

type AppFullInfo struct {
	AppInfo
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func (i *AppInfo) Token(ts int64) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s:%v:%s", i.AppId, i.AppKey, ts, getSecretString(gconv.String(ts), i.AppSecret))),
	)
}

func (i *AppInfo) FromToken(token string) (ts, secret string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 4 {
		err = errors.New("invalid token parts")
		return
	}

	i.AppId = parts[0]
	i.AppKey = parts[1]
	ts = parts[2]
	secret = parts[3]
	return
}

func (i *AppInfo) Check(appKey, ts, secret string) (err error) {
	if i.AppKey != appKey {
		err = errors.New("invalid app key")
		return
	}
	if getSecretString(ts, i.AppSecret) != secret {
		err = errors.New("invalid app secret")
		return
	}
	return
}

func getSecretString(ts, secret string) string {
	return gmd5.MustEncryptString(fmt.Sprintf("%v:%v", gmd5.MustEncrypt(secret), ts))
}
