package gokafka

type Config struct {
	User     string   `json:"user" yaml:"User"`
	Name     string   `json:"name" yaml:"Name"`
	Password string   `json:"password" yaml:"Password"`
	Addrs    []string `json:"addrs" yaml:"Addrs"`
	Timeout  int      `json:"timeout" yaml:"Timeout"`
}
