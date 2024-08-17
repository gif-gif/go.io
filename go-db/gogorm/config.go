package gogorm

import (
	"gorm.io/gorm"
	"time"
)

// 数据初始化属性这里扩展
type GormDbConfig struct {
	//AutoPing     bool
	Config       *gorm.Config
	MaxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           // <= 0 means unlimited
	MaxLifetime  time.Duration // maximum amount of time a connection may be reused
}

type Config struct {
	DataSource   string `yaml:"DataSource,optional" json:"dataSource,optional"`
	MaxIdleCount int    `yaml:"MaxIdleCount,optional" json:"maxIdleCount,optional"` // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int    `yaml:"MaxOpen,optional" json:"maxOpen,optional"`           // <= 0 means unlimited
	MaxLifetime  int    `yaml:"MaxLifetime,optional" json:"maxLifetime,optional"`   // maximum amount of time a connection may be reused
}
