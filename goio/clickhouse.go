package goio

import (
	goclickhouse "github.com/gif-gif/go.io/go-db/go-clickhouse"
)

func CK(names ...string) *goclickhouse.GoClickHouse {
	return goclickhouse.GetClient(names...)
}
