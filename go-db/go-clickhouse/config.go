package goclickhouse

type Config struct {
	Name string `yaml:"Name" json:"name,optional"`
	//Driver             string `yaml:"Driver" json:"driver,optional"`
	Addr               []string `yaml:"Addr" json:"addr,optional"`
	User               string   `yaml:"User" json:"user,optional"`
	Password           string   `yaml:"Password" json:"password,optional"`
	Database           string   `yaml:"Database" json:"database,optional"`
	DialTimeout        int32    `yaml:"DialTimeout" json:"dialTimeout,optional"`               // default 30 second
	MaxIdleConn        int      `yaml:"MaxIdleConn" json:"maxIdleConn,optional"`               // default 5 second
	MaxOpenConn        int      `yaml:"MaxOpenConn" json:"maxOpenConn,optional"`               // default 10 second
	MaxExecutionTime   int      `yaml:"MaxExecutionTime" json:"maxExecutionTime,optional"`     // default 60 second
	ConnMaxLifetime    int      `yaml:"ConnMaxLifetime" json:"connMaxLifetime,optional"`       //seconds 60 * 60（1hour）
	TLS                bool     `yaml:"TLS" json:"tls,optional"`                               // tls true 时都会启用https 否则http
	InsecureSkipVerify bool     `yaml:"InsecureSkipVerify" json:"insecureSkipVerify,optional"` // tls true 才会生效
	AutoPing           bool     `yaml:"AutoPing" json:"autoPing,optional"`
	Debug              bool     `yaml:"Debug" json:"debug,optional"`
}
