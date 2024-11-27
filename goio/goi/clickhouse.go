package goi

import (
	"database/sql"
	goclickhouse "github.com/gif-gif/go.io/go-db/go-clickhouse"
)

func CkClient(names ...string) *goclickhouse.GoClickHouse {
	return goclickhouse.GetClient(names...)
}

func CK(names ...string) *sql.DB {
	client := CkClient(names...)
	if client == nil {
		return nil
	}
	return client.DB()
}
