package user

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
)

func TestPlugin(t *testing.T) {
	s := g.Server()
	s.SetPort(8000)

	s.Plugin(NewPlugin())
	s.Run()
}
