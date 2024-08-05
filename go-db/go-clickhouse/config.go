package goclickhouse

type Config struct {
	Driver             string `yaml:"Driver,optional" json:"driver,optional"`
	Addr               string `yaml:"Addr,optional" json:"addr,optional"`
	User               string `yaml:"User,optional" json:"user,optional"`
	Password           string `yaml:"Password,optional" json:"password,optional"`
	Database           string `yaml:"Database,optional" json:"database,optional"`
	DialTimeout        int32  `yaml:"DialTimeout,optional" json:"dialTimeout,optional"`               // default 30 second
	MaxIdleConn        int    `yaml:"MaxIdleConn,optional" json:"maxIdleConn,optional"`               // default 5 second
	MaxOpenConn        int    `yaml:"MaxOpenConn,optional" json:"maxOpenConn,optional"`               // default 10 second
	MaxExecutionTime   int    `yaml:"MaxExecutionTime,optional" json:"maxExecutionTime,optional"`     // default 60 second
	ConnMaxLifetime    int    `yaml:"ConnMaxLifetime,optional" json:"connMaxLifetime,optional"`       //seconds 60 * 60（1hour）
	TLS                bool   `yaml:"TLS,optional" json:"tls,optional"`                               // tls true 时都会启用https 否则http
	InsecureSkipVerify bool   `yaml:"InsecureSkipVerify,optional" json:"insecureSkipVerify,optional"` // tls true 才会生效
	AutoPing           bool   `yaml:"AutoPing,optional" json:"autoPing,optional"`
	Debug              bool   `yaml:"Debug,optional" json:"debug,optional"`
}
