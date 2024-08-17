package goip

type Config struct {
	Mmdb          string `yaml:"Mmdb" json:"mmdb,optional"`
	Ip2locationDB string `yaml:"Ip2locationDB" json:"ip2locationDB,optional"`
}
