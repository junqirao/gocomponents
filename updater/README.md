# Updater

execute function or embed sql script.

## How dose it work

### standalone

1. fetch updater execute record
2. compare and execute: skip failure or cannot retry
3. create or update record

### cluster

every node:

1. try get distribute locker
2. fetch updater execute record
3. compare and execute: skip failure or cannot retry
4. create or update record
5. release lock

## Usage

```golang
package main

import (
	"context"
	"embed"

	"github.com/junqirao/gocomponents/dkv"
	"github.com/junqirao/gocomponents/updater"
	"github.com/gogf/gf/v2/frame/g"
)

//go:embed somewhere
var EmbeddedFS embed.FS

func main() {
	ctx := context.Background()
	db, err := dkv.NewDB(ctx)
	if err != nil {
		return
	}

	// define functions
	// generate sql script function by read embed.FS sql/*.sql
	// create function named as sql_{script_name}
	fis := updater.SQLFuncFromEmbedFS(ctx, g.DB(), EmbeddedFS)
	// name  : unique name with type
	// fn    : function
	// must  : if must and error caused by this function updater will terminate
	// retry : function can retry
	// typ   : specific type, default 0=raw
	fis = append(fis, updater.NewFunc("my_update_func", func(ctx context.Context) (err error) {
		// do something
		return nil
	}, true, true))

	// use distribute kv database as record store database
	adaptor := updater.NewKVDatabaseAdaptor(db)
	// or use mysql as record store database
	// adaptor :=updater.NewMysqlAdaptor(ctx, g.DB(), g.Log())

	// execute functions
	// standalone, no distribute lock
	// err = updater.Update2Latest(ctx, updater.NewKVDatabaseAdaptor(db), fis...)
	// cluster with distribute lock
	err = updater.ConcurrencyUpdate2Latest(ctx, adaptor, db, fis...)
	if err != nil {
		return
	}
}

```