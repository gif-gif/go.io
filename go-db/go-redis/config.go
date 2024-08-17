package goredis

type Config struct {
	Name     string `yaml:"Name" json:"name,optional"`
	Addr     string `yaml:"Addr" json:"addr,optional"`
	Password string `yaml:"Password" json:"password,optional"`
	DB       int    `yaml:"DB" json:"db,optional"`
	Prefix   string `yaml:"Prefix" json:"prefix,optional"`
	AutoPing bool   `yaml:"AutoPing" json:"autoPing,optional"`
}
