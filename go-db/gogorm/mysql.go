package gogorm

import (
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化MySQL OR TiDB(TiDB is compatible with MySQL protocol.)
func InitMysql(dataSource string, config GormDbConfig) (*GormDB, error) {
	golog.Info("init Mysql")
	if config.Config == nil {
		config.Config = &gorm.Config{QueryFields: true}
	}
	//config.Config.DisableAutomaticPing = !config.AutoPing
	db, err := gorm.Open(mysql.Open(dataSource), config.Config)
	if err != nil {
		return nil, err
	}

	m := &GormDB{
		DB: db,
	}

	err = m.Init(&config)
	if err != nil {
		return nil, err
	}

	return m, nil
}
