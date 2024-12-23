package controller

import (
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/user/logic"
)

var Security = &security{}

type security struct {
}

func (security) GetTransportPublicKeyPem(r *ghttp.Request) {
	pem, err := logic.Security.GetTransportPublicKeyPem(r.Context())
	if err != nil {
		response.Error(r, err)
		return
	}
	response.Success(r, gbase64.EncodeString(pem))
}
