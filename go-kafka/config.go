package gokafka

type Config struct {
	User     string   `json:"user" yaml:"user"`
	Password string   `json:"password" yaml:"password"`
	Addrs    []string `json:"addrs" yaml:"addrs"`
	Timeout  int      `json:"timeout" yaml:"timeout"`
}
