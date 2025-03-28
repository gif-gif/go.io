package gogorm

import (
	"errors"
	"github.com/gif-gif/go.io/goio"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type GoGorm struct {
	DB     *gorm.DB
	Config *Config
}

func New(config *Config) (*GoGorm, error) {
	dlt, err := createDialector(config)
	if err != nil {
		return nil, err
	}

	if config.GormConfig == nil {
		config.GormConfig = createDefaultConfig(config)
	}

	db, err := gorm.Open(dlt, config.GormConfig)
	if err != nil {
		return nil, err
	}

	m := &GoGorm{
		DB:     db,
		Config: config,
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if sqlDB != nil {
		if config.MaxIdleCount == 0 {
			config.MaxIdleCount = 10
		}
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(config.MaxIdleCount)
		if config.MaxOpen == 0 {
			config.MaxOpen = 100
		}
		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(config.MaxOpen)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		if config.MaxLifetime == 0 {
			sqlDB.SetConnMaxLifetime(time.Hour)
		} else {
			sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Minute)
		}
	}

	return m, nil
}

func createDialector(config *Config) (gorm.Dialector, error) {
	switch config.DBType {
	case DATABASE_MYSQL:
	case DATABASE_TIDB:
		return mysql.Open(config.DataSource), nil
	case DATABASE_STARROCKS:
		return mysql.Open(config.DataSource), nil
	case DATABASE_POSTGRESQL:
		return postgres.New(postgres.Config{
			DSN:                  config.DataSource,
			PreferSimpleProtocol: false, // enable implicit prepared statement usage
		}), nil
	case DATABASE_SQLITE:
		return sqlite.Open(config.DataSource), nil
	case DATABASE_SQLSERVER:
		return sqlserver.Open(config.DataSource), nil
	case DATABASE_CLICKHOUSE:
		return clickhouse.Open(config.DataSource), nil
	default:
		return mysql.Open(config.DataSource), nil
	}

	return nil, errors.New("unsupported database type")
}

func createDefaultConfig(config *Config) *gorm.Config {
	logLevel := logger.Info
	if goio.IsPro() {
		logLevel = logger.Silent
	}
	switch config.DBType {
	case DATABASE_MYSQL:
	case DATABASE_TIDB:
		return &gorm.Config{Logger: logger.Default.LogMode(logLevel)}
	case DATABASE_STARROCKS:
		return &gorm.Config{
			// 使用较详细的日志级别，生产环境可调整
			Logger: logger.Default.LogMode(logLevel),
			// 关闭自动复数化命名，StarRocks对此支持有限
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			// 禁用外键约束，StarRocks不支持外键
			DisableForeignKeyConstraintWhenMigrating: true,
		}
	case DATABASE_POSTGRESQL:
		return &gorm.Config{Logger: logger.Default.LogMode(logLevel)}
	case DATABASE_SQLITE:
		return &gorm.Config{Logger: logger.Default.LogMode(logLevel)}
	case DATABASE_SQLSERVER:
		return &gorm.Config{Logger: logger.Default.LogMode(logLevel)}
	case DATABASE_CLICKHOUSE:
		return &gorm.Config{Logger: logger.Default.LogMode(logLevel)}
	default:
		return &gorm.Config{Logger: logger.Default.LogMode(logLevel)}
	}

	return nil
}
