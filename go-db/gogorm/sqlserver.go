package gogorm

import (
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// 初始化
// dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
func InitSqlServer(dsn string, config GormDbConfig) (*GormDB, error) {
	golog.Info("init PostgreSql")
	if config.Config == nil {
		config.Config = &gorm.Config{QueryFields: true}
	}
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	s := &GormDB{
		DB: db,
	}

	err = s.Init(&config)
	if err != nil {
		return nil, err
	}
	return s, nil
}
