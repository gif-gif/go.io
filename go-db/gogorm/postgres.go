package gogorm

import (
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 初始化
// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
func InitPostgreSql(dsn string, config GormDbConfig) (*GormDB, error) {
	golog.Info("init PostgreSql")
	if config.Config == nil {
		config.Config = &gorm.Config{QueryFields: true}
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false, // enable implicit prepared statement usage
	}), &gorm.Config{})

	//db, err := gorm.Open(postgres.Open(dsn), config.Config)
	//if err != nil {
	//	return nil, err
	//}

	s := &GormDB{
		DB: db,
	}

	err = s.Init(&config)
	if err != nil {
		return nil, err
	}
	return s, nil
}
