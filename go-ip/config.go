package goip

type Config struct {
	Name          string `yaml:"Name" json:"name,optional"`
	Mmdb          string `yaml:"Mmdb" json:"mmdb,optional"`
	Ip2locationDB string `yaml:"Ip2locationDB" json:"ip2locationDB,optional"`
	IpServiceUrl  string `yaml:"IpServiceUrl" json:"ipServiceUrl,optional"`
}
