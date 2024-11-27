package goi

import (
	"github.com/gif-gif/go.io/go-db/gogorm"
	"gorm.io/gorm"
)

func DbClient(names ...string) *gogorm.GoGorm {
	return gogorm.GetClient(names...)
}

func DB(names ...string) *gorm.DB {
	client := DbClient(names...)
	if client == nil {
		return nil
	}
	return client.DB
}
