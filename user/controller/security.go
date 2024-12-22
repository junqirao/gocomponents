package controller

import (
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/user/logic"
)

func GetPublicKeyPem(r *ghttp.Request) {
	pem, err := logic.GetSecurityPublicKeyPem(r.Context())
	if err != nil {
		response.Error(r, err)
		return
	}
	response.Success(r, gbase64.EncodeString(pem))
}
