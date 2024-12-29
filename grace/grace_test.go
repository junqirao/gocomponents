package grace

import (
	"context"
	"testing"
)

func TestGracefulExit(t *testing.T) {
	ctx := context.TODO()
	Register(ctx, "test", func() {
		t.Log("test")
	})
	GracefulExit(ctx)
}
