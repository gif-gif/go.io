package gogorm

type Config struct {
	Name         string `yaml:"Name" json:"name,optional"`
	DBType       string `yaml:"DBType" json:"dbType,optional"` // default mysql
	DataSource   string `yaml:"DataSource" json:"dataSource,optional"`
	MaxIdleCount int    `yaml:"MaxIdleCount" json:"maxIdleCount,optional"` // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int    `yaml:"MaxOpen" json:"maxOpen,optional"`           // <= 0 means unlimited
	MaxLifetime  int    `yaml:"MaxLifetime" json:"maxLifetime,optional"`   // maximum amount of time a connection may be reused minutes
}

const (
	DATABASE_MYSQL      = "mysql"
	DATABASE_POSTGRESQL = "postgres"
	DATABASE_SQLITE     = "sqlite"
	DATABASE_SQLSERVER  = "sqlserver"
	DATABASE_CLICKHOUSE = "clickhouse"
	DATABASE_TIDB       = "tidb"
	DATABASE_STARROCKS  = "starrocks"
)
