package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/user/logic"
	"github.com/junqirao/gocomponents/user/model"
)

// CreateUser creates a new user
func CreateUser(r *ghttp.Request) {
	params := new(model.CreateUserReq)
	if err := r.Parse(params); err != nil {
		response.Error(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	user, err := logic.CreateUser(r.Context(), params)
	if err != nil {
		response.Error(r, err)
		return
	}

	response.Success(r, user)
}

// CheckUsernameExists checks whether the username exists
func CheckUsernameExists(r *ghttp.Request) {
	params := new(model.CheckUsernameReq)
	if err := r.Parse(params); err != nil {
		response.Error(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	err := logic.UserExists(r.Context(), params.Username)
	if err != nil {
		response.Error(r, err)
		return
	}

	response.Success(r)
}
