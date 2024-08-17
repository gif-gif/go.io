package gokafka

type Config struct {
	User     string   `json:"user,optional" yaml:"User,omitempty"`
	Name     string   `json:"name,optional"  yaml:"Name"`
	Password string   `json:"password,optional"  yaml:"Password"`
	Addrs    []string `json:"addrs,optional"  yaml:"Addrs"`
	Timeout  int      `json:"timeout,optional"  yaml:"Timeout"`
}
