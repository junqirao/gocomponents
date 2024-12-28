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

func NewMutex(ctx context.Context, database Database, name string) (mu Mutex, err error) {
	locker, err := database.Locker(ctx, fmt.Sprintf("lock_%s", name))
	if err != nil {
		return
	}
	return Mutex{key: fmt.Sprintf("lock_%s", name), Locker: locker}, nil
}
