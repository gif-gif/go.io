package godb

import (
	"gorm.io/gorm"
	"time"
)

// 数据初始化属性这里扩展
type GoDbConfig struct {
	Config       *gorm.Config
	MaxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           // <= 0 means unlimited
	MaxLifetime  time.Duration // maximum amount of time a connection may be reused
}
