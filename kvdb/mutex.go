package kvdb

import (
	"context"
	"fmt"
	"sync"
)

type Mutex struct {
	key string
	sync.Locker
}

func NewMutex(ctx context.Context, name string, database ...Database) (mu Mutex, err error) {
	var db = Raw
	if len(database) > 0 {
		db = database[0]
	}

	locker, err := db.Locker(ctx, fmt.Sprintf("lock_%s", name))
	if err != nil {
		return
	}
	return Mutex{key: fmt.Sprintf("lock_%s", name), Locker: locker}, nil
}
