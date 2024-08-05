package goip

type Config struct {
	Mmdb          string `yaml:"Mmdb,optional" json:"mmdb,optional"`
	Ip2locationDB string `yaml:"Ip2locationDB,optional" json:"ip2locationDB,optional"`
}
