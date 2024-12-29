package audit

import (
	"context"
)

const (
	EventShutdown = "shutdown"
	EventStartUp  = "startup"
)

const (
	ModuleSystem = "system"
)

func Sys(ctx context.Context, event string, content ...interface{}) {
	Logger.Ctx(ctx).Module(ModuleSystem).Log(event, content)
}
