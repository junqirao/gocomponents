package updater

import (
	"context"
	"fmt"
	"testing"

	emb "github.com/junqirao/gocomponents/updater/embed"
)

func TestSQLFuncFromEmbedFS(t *testing.T) {
	fis := SQLFuncFromEmbedFS(context.Background(), nil, emb.EmbeddedFS)
	for _, fi := range fis {
		t.Log(fmt.Sprintf("name=%s, type=%v", fi.Name, fi.Type))
	}
}
