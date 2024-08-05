package goredis

type Config struct {
	Name     string `yaml:"Name,optional" json:"name,optional"`
	Addr     string `yaml:"Addr,optional" json:"addr,optional"`
	Password string `yaml:"Password,optional" json:"password,optional"`
	DB       int    `yaml:"DB,optional" json:"db,optional"`
	Prefix   string `yaml:"Prefix,optional" json:"prefix,optional"`
	AutoPing bool   `yaml:"AutoPing,optional" json:"autoPing,optional"`
}
