package gogorm

import (
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 初始化sqlite3
//
// Example: dbFilePath : /user/db/sqlite3.db
func InitSqlite3(dbFilePath string, config GormDbConfig) (*GormDB, error) {
	golog.Info("init Sqlite3")
	if config.Config == nil {
		config.Config = &gorm.Config{QueryFields: true}
	}
	//config.Config.DisableAutomaticPing = !config.AutoPing
	db, err := gorm.Open(sqlite.Open(dbFilePath), config.Config)
	if err != nil {
		return nil, err
	}

	s := &GormDB{
		DB: db,
	}

	err = s.Init(&config)
	if err != nil {
		return nil, err
	}
	return s, nil
}
