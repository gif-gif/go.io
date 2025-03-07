package gokafka

import goredis "github.com/gif-gif/go.io/go-db/go-redis"

type Config struct {
	User              string   `json:"user,optional" yaml:"User"`
	Name              string   `json:"name,optional"  yaml:"Name"`
	Password          string   `json:"password,optional"  yaml:"Password"`
	Addrs             []string `json:"addrs,optional"  yaml:"Addrs"`
	Timeout           int      `json:"timeout,optional"  yaml:"Timeout"`
	HeartbeatInterval int      `json:"heartbeatInterval,optional" yaml:"HeartbeatInterval"`
	SessionTimeout    int      `json:"sessionTimeout,optional" yaml:"SessionTimeout"`
	RebalanceTimeout  int      `json:"rebalanceTimeout,optional" yaml:"RebalanceTimeout"`
	OffsetNewest      bool     `json:"offsetNewest,optional" yaml:"OffsetNewest"`
	Version           string   `json:"version,optional" yaml:"Version"`

	GroupId     string         `json:"groupId,optional" yaml:"GroupId"`
	RedisConfig goredis.Config `json:"redisConfig,optional" yaml:"RedisConfig"`
}
