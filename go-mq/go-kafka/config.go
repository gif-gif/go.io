package gokafka

type Config struct {
	User              string   `json:"user,optional" yaml:"User"`
	Name              string   `json:"name,optional"  yaml:"Name"`
	Password          string   `json:"password,optional"  yaml:"Password"`
	Addrs             []string `json:"addrs,optional"  yaml:"Addrs"`
	Timeout           int      `json:"timeout,optional"  yaml:"Timeout"`
	HeartbeatInterval int      `json:"heartbeatInterval" yaml:"HeartbeatInterval"`
	SessionTimeout    int      `json:"sessionTimeout" yaml:"SessionTimeout"`
	RebalanceTimeout  int      `json:"rebalanceTimeout" yaml:"RebalanceTimeout"`
	OffsetNewest      bool     `json:"offsetNewest" yaml:"OffsetNewest"`
}
