package goclickhouse

type Config struct {
	Driver             string `yaml:"Driver" json:"Driver"`
	Addr               string `yaml:"Addr" json:"Addr"`
	User               string `yaml:"User" json:"User"`
	Password           string `yaml:"Password" json:"Password"`
	Database           string `yaml:"Database" json:"Database"`
	DialTimeout        int32  `yaml:"DialTimeout" json:"DialTimeout"`               // default 30 second
	MaxIdleConn        int    `yaml:"MaxIdleConn" json:"MaxIdleConn"`               // default 5 second
	MaxOpenConn        int    `yaml:"MaxOpenConn" json:"MaxOpenConn"`               // default 10 second
	MaxExecutionTime   int    `yaml:"max_execution_time" json:"MaxExecutionTime"`   // default 60 second
	ConnMaxLifetime    int    `yaml:"ConnMaxLifetime" json:"ConnMaxLifetime"`       //seconds 60 * 60（1hour）
	Tls                bool   `yaml:"Tls" json:"Tls"`                               // tls true 时都会启用https 否则http
	InsecureSkipVerify bool   `yaml:"InsecureSkipVerify" json:"InsecureSkipVerify"` // tls true 才会生效
	AutoPing           bool   `yaml:"AutoPing" json:"AutoPing"`
	Debug              bool   `yaml:"Debug" json:"Debug"`
}
