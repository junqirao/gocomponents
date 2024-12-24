package updater

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SQLFuncFromEmbedFS read sql/*.sql from embed.FS,
// convert *.sql to FuncInfo named of sql_*
func SQLFuncFromEmbedFS(ctx context.Context, db gdb.DB, fs embed.FS, inline ...bool) (fis []*FuncInfo) {
	fis = make([]*FuncInfo, 0)
	dir, err := fs.ReadDir("sql")
	if err != nil {
		return
	}
	for _, entry := range dir {
		if !entry.IsDir() &&
			!strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		var content []byte
		content, err = readContentFromEntry("sql/", fs, entry)
		if err != nil {
			g.Log().Warningf(ctx, "read sql/%s error: %v", entry.Name(), err)
			continue
		}
		if len(inline) > 0 && inline[0] {
			content = removeLineBreak(content)
		}
		fis = append(fis,
			NewFunc(
				fmt.Sprintf("sql_%s", strings.TrimSuffix(entry.Name(), ".sql")),
				generateSqlFunc(content, db),
				NewFuncConfig().Retry().Must(), FuncTypeSql,
			),
		)
	}
	return
}

func removeLineBreak(content []byte) []byte {
	s := string(content)
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return []byte(s)
}

func generateSqlFunc(content []byte, db gdb.DB) FN {
	return func(ctx context.Context) (err error) {
		if db == nil {
			return
		}
		_, err = db.Exec(ctx, string(content))
		return
	}
}

func readContentFromEntry(prefix string, fs embed.FS, entry fs.DirEntry) (bs []byte, err error) {
	info, err := entry.Info()
	if err != nil {
		return
	}
	return fs.ReadFile(prefix + info.Name())
}
