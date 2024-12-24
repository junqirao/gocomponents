package autoconfig

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/auth"
)

func init() {
	ctx := context.Background()
	v, err := g.Cfg().Get(ctx, "auth")
	if err != nil {
		g.Log().Errorf(ctx, "failed to get auth config: %v", err)
		return
	}

	cfg := new(auth.Config)
	if err = v.Scan(&cfg); err != nil {
		g.Log().Errorf(ctx, "failed to read auth config: %v", err)
		return
	}

	auth.SetConfig(cfg)
}
