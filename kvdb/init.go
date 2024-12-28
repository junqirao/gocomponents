package kvdb

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	var (
		ctx    = gctx.GetInitCtx()
		config = Config{}
	)

	v, err := g.Cfg().Get(ctx, "g.Cfg().Get()")
	if err != nil {
		g.Log().Errorf(ctx, "kvdb failed to get config: %v", err)
		return
	}
	if err = v.Struct(&config); err != nil {
		g.Log().Errorf(ctx, "kvdb failed to parse config: %v", err)
		return
	}

	config.check()
	Raw, err = newEtcd(ctx, config)
	// create Storages instance
	Storages = newStorages(ctx, config, Raw)
	return
}
