package gokafka

type Config struct {
	User     string   `json:"user,optional" yaml:"User,optional"`
	Name     string   `json:"name,optional" yaml:"Name,optional"`
	Password string   `json:"password,optional" yaml:"Password,optional"`
	Addrs    []string `json:"addrs,optional" yaml:"Addrs,optional"`
	Timeout  int      `json:"timeout,optional" yaml:"Timeout,optional"`
}
