package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	uuid "github.com/satori/go.uuid"

	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/types"
	"github.com/junqirao/gocomponents/user/dao"
	"github.com/junqirao/gocomponents/user/model"
	"github.com/junqirao/gocomponents/user/model/entity"
)

func UserExists(ctx context.Context, username string) (err error) {
	cnt, err := dao.CUser.Ctx(ctx).Count(g.Map{
		dao.CUser.Columns().Username: username,
	})
	if err != nil {
		return
	}

	if cnt > 0 {
		err = response.CodeConflict.WithDetail(fmt.Sprintf("username [%s] already exists", username))
	}
	return
}

func CreateUser(ctx context.Context, input *model.CreateUserReq) (user *types.User, err error) {
	if err = UserExists(ctx, input.Username); err != nil {
		return
	}

	eu := entity.CUser{
		Id:            uuid.NewV4().String(),
		Username:      input.Username,
		Password:      input.Password,
		CreatedAt:     gtime.Now(),
		UpdatedAt:     nil,
		Administrator: 0,
		Source:        types.UserSourceInternal,
		Status:        types.UserStatusActive.Int(),
		Extra:         "{}",
	}

	if _, err = dao.CUser.Ctx(ctx).Insert(eu); err != nil {
		return
	}

	user = &types.User{
		Id:            eu.Id,
		Username:      eu.Username,
		CreatedAt:     eu.CreatedAt.UnixMilli(),
		Administrator: eu.Administrator == 1,
		Source:        eu.Source,
		Status:        types.UserStatus(eu.Status),
		Extra:         make(map[string]any),
	}
	if eu.Extra != "" {
		_ = json.Unmarshal([]byte(eu.Extra), &user.Extra)
	}
	return
}

func CreateAdminIfNotExists(ctx context.Context) (err error) {
	count, err := dao.CUser.Ctx(ctx).Count(g.Map{
		dao.CUser.Columns().Administrator: 1,
	})
	if err != nil {
		return
	}
	if count > 0 {
		return
	}
	pwd := grand.S(16)
	username := "admin"

	if err := UserExists(ctx, username); err != nil {
		username = fmt.Sprintf("admin_%v", grand.S(5))
	}

	_, err = dao.CUser.Ctx(ctx).Insert(entity.CUser{
		Id:            uuid.NewV4().String(),
		Username:      username,
		Password:      pwd,
		CreatedAt:     gtime.Now(),
		UpdatedAt:     nil,
		Administrator: 1,
		Source:        types.UserSourceInternal,
		Status:        types.UserStatusActive.Int(),
		Extra:         "{}",
	})
	if err == nil {
		g.Log().Infof(ctx, `

-------------------------
admin account auto created, please change your password after login as soon as possible.
username: %s
password: %s
-------------------------

`, username, pwd)
	}
	return
}
