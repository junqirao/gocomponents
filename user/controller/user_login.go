package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/user/logic"
	"github.com/junqirao/gocomponents/user/model"
)

func (user) Login(r *ghttp.Request) {
	params := new(model.UserLoginReq)
	if err := r.Parse(params); err != nil {
		response.Error(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	u, err := logic.User.Login(r.Context(), params)
	if err != nil {
		response.Error(r, err)
		return
	}

	response.Success(r, u)
}
