package gogorm

import (
	"errors"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
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
	db, err := gorm.Open(dlt, &gorm.Config{QueryFields: true})
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
