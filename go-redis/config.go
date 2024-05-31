package goredis

type Config struct {
	Name     string `yaml:"name" json:"name"`
	Addr     string `yaml:"addr" json:"addr"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"DB"`
	Prefix   string `yaml:"prefix" json:"prefix"`
	AutoPing bool   `yaml:"auto_ping" json:"autoPing"`
}
