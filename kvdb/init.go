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
	err := g.Cfg().MustGet(ctx, "kvdb").Struct(&config)
	if err != nil {
		g.Log().Errorf(ctx, "kvdb failed to parse config: %v", err)
		return
	}

	config.check()
	Raw, err = newEtcd(ctx, config)
	if err != nil {
		panic(err)
	}
	// create Storages instance
	Storages = newStorages(ctx, config, Raw)
	return
}
