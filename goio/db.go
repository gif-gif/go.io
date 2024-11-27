package goio

import (
	"github.com/gif-gif/go.io/go-db/gogorm"
	"gorm.io/gorm"
)

func GODB(names ...string) *gogorm.GoGorm {
	return gogorm.GetClient(names...)
}

func DB(names ...string) *gorm.DB {
	client := GODB(names...)
	if client == nil {
		return nil
	}
	return client.DB
}
