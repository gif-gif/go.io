package goi

import (
	gominio "github.com/gif-gif/go.io/go-oss/go-minio"
)

func GoMinio(names ...string) *gominio.Uploader {
	return gominio.GetClient(names...)
}
