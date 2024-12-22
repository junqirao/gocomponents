package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	uuid "github.com/satori/go.uuid"

	"github.com/junqirao/gocomponents/jwt"
	"github.com/junqirao/gocomponents/response"
	"github.com/junqirao/gocomponents/security"
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

func UserLogin(ctx context.Context, input *model.UserLoginReq) (user *types.UserInfoWithToken, err error) {
	u := new(entity.CUser)
	eu, err := dao.CUser.Ctx(ctx).Where(dao.CUser.Columns().Username, input.Username).One()
	if err != nil {
		return
	}
	if err = eu.Struct(u); errors.Is(gerror.Cause(err), sql.ErrNoRows) {
		err = response.CodeNotFound.WithMessage("invalid user or password")
		return
	}

	sp, err := security.GetProvider(ctx, securityProviderNameTransport)
	if err != nil {
		return
	}

	var dbPwd = u.Password
	if decrypt, err := sp.Decrypt(u.Password); err == nil {
		// supports raw storage
		dbPwd = decrypt
	}

	iptPwd, err := sp.Decrypt(input.Password)
	if err != nil {
		return
	}

	if dbPwd != iptPwd {
		err = response.CodeUnauthorized.WithDetail("invalid user or password")
		return
	}

	accessToken, err := jwt.GenerateToken(&jwt.GenerateTokenRequest{
		UserId:    u.Id,
		From:      input.From,
		UserName:  u.Username,
		Subject:   fmt.Sprintf("c_user_%s", u.Username),
		Issuer:    "c_user",
		ExpiredIn: accessTokenExpireTime,
		Key:       tokenKey,
	})
	if err != nil {
		return
	}

	user = &types.UserInfoWithToken{
		User: types.User{
			Id:            u.Id,
			Username:      u.Username,
			CreatedAt:     u.CreatedAt.UnixMilli(),
			Administrator: u.Administrator == 1,
			Source:        u.Source,
			Status:        types.UserStatus(u.Status),
			Extra:         make(map[string]any),
		},
		AccessToken: "Bearer " + accessToken,
	}
	if u.Extra != "" {
		_ = json.Unmarshal([]byte(u.Extra), &user.Extra)
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
