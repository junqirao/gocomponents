package controller

import (
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/security"
)

func GetPublicKeyPem(r *ghttp.Request) {
	response.Success(r, gbase64.EncodeString(security.GetPublicKeyPem()))
}
