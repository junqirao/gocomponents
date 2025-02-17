package kvdb

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func Init(ctx context.Context) (err error) {
	var (
		config = Config{}
	)
	if err = g.Cfg().MustGet(ctx, "kvdb.database").Struct(&config); err != nil {
		g.Log().Errorf(ctx, "kvdb failed to parse config: %v", err)
		return
	}

	Raw, err = newEtcd(ctx, config)
	return
}

func InitStorage(ctx context.Context, db ...Database) (err error) {
	var (
		config = StorageConfig{}
	)
	if err = g.Cfg().MustGet(ctx, "kvdb.storage").Struct(&config); err != nil {
		g.Log().Errorf(ctx, "kvdb failed to parse config: %v", err)
		return
	}
	if config.Separator == "" {
		config.Separator = "/"
	}
	database := Raw
	if len(db) > 0 && db[0] != nil {
		database = db[0]
	}
	// create Storages instance
	Storages = newStorages(ctx, config, database)
	return
}
