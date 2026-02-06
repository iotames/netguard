package sql

import (
	"embed"
)

//go:embed *.sql
var sqlFS embed.FS

func GetSqlFs() embed.FS {
	return sqlFS
}
