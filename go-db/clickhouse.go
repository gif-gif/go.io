package godb

import (
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

// 初始化CLICKHOUSE
func InitClickHouse(dataSource string, config GoDbConfig) (*GoDB, error) {
	golog.Info("init ClickHouse")
	if config.Config == nil {
		config.Config = &gorm.Config{QueryFields: true}
	}
	//config.Config.DisableAutomaticPing = !config.AutoPing
	db, err := gorm.Open(clickhouse.Open(dataSource), config.Config)
	if err != nil {
		return nil, err
	}

	m := &GoDB{
		DB: db,
	}

	err = m.Init(&config)
	if err != nil {
		return nil, err
	}

	return m, nil
}
