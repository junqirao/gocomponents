package embed

import (
	"embed"
)

//go:embed sql
var EmbeddedFS embed.FS

func GetInitSqlContent() string {
	file, err := EmbeddedFS.ReadFile("sql/init.sql")
	if err != nil {
		return ""
	}
	return string(file)
}
