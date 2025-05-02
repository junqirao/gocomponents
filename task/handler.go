package task

import (
	"sync"
)

var (
	handlers = sync.Map{} // name:Handler
)

func RegisterHandler(name string, handler Handler) {
	handlers.Store(name, handler)
}

func getHandler(typ string) (h Handler, err error) {
	v, ok := handlers.Load(typ)
	if !ok {
		err = ErrTaskTypeNotFound
		return
	}
	h, ok = v.(Handler)
	if !ok {
		err = ErrTaskTypeNotFound
		return
	}
	return
}
