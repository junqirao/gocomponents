package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/user/logic"
	"github.com/junqirao/gocomponents/user/model"
)

var User = &user{}

type user struct {
}

// Create creates a new user
func (user) Create(r *ghttp.Request) {
	params := new(model.CreateUserReq)
	if err := r.Parse(params); err != nil {
		response.Error(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	u, err := logic.User.Create(r.Context(), params)
	if err != nil {
		response.Error(r, err)
		return
	}

	response.Success(r, u)
}

// CheckUsernameExists checks whether the username exists
func (user) CheckUsernameExists(r *ghttp.Request) {
	params := new(model.CheckUsernameReq)
	if err := r.Parse(params); err != nil {
		response.Error(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	err := logic.User.Exists(r.Context(), params.Username)
	if err != nil {
		response.Error(r, err)
		return
	}

	response.Success(r)
}
