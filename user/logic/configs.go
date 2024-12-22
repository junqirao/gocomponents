package logic

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	accessTokenExpireTime = tryGetConfigWithDefault("user.access_token_expire_time", 60*60*2).Int64()
	tokenKey              = tryGetConfigWithDefault("user.token_key", "go_component_user_module").Bytes()
)

func tryGetConfigWithDefault(pattern string, def any) *g.Var {
	v, err := g.Cfg().Get(gctx.GetInitCtx(), pattern)
	if err != nil || v.IsNil() || v.IsEmpty() {
		return g.NewVar(def)
	}
	return v
}
