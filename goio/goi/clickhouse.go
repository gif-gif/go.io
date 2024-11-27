package goi

import (
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	goclickhouse "github.com/gif-gif/go.io/go-db/go-clickhouse"
)

func CkClient(names ...string) *goclickhouse.GoClickHouse {
	return goclickhouse.GetClient(names...)
}

func CK(names ...string) driver.Conn {
	client := CkClient(names...)
	if client == nil {
		return nil
	}
	return client.Conn()
}

func CKDB(names ...string) *sql.DB {
	client := CkClient(names...)
	if client == nil {
		return nil
	}
	return client.DB()
}
