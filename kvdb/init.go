package kvdb

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func Init(ctx context.Context) (err error) {
	databaseOnceInit.Do(func() {
		var (
			config = Config{}
		)
		if err = g.Cfg().MustGet(ctx, "kvdb.database").Struct(&config); err != nil {
			g.Log().Errorf(ctx, "kvdb failed to parse config: %v", err)
			return
		}

		Raw, err = newEtcd(ctx, config)
	})
	return
}

func MustGetDatabase(ctx context.Context) Database {
	if Raw == nil {
		if err := Init(ctx); err != nil {
			panic(err)
		}
	}
	return Raw
}

func InitStorage(ctx context.Context, db ...Database) (err error) {
	storageOnceInit.Do(func() {
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
		database := MustGetDatabase(ctx)
		if len(db) > 0 && db[0] != nil {
			database = db[0]
		}
		// create Storages instance
		Storages = newStorages(ctx, config, database)
	})
	return
}
