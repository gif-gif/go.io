package gokafka

import goredis "github.com/gif-gif/go.io/go-db/go-redis"

type Config struct {
	User              string   `json:"user,optional" yaml:"User"`
	Name              string   `json:"name,optional"  yaml:"Name"`
	Password          string   `json:"password,optional"  yaml:"Password"`
	Addrs             []string `json:"addrs,optional"  yaml:"Addrs"`
	Timeout           int      `json:"timeout,optional"  yaml:"Timeout"`                    // 单位：秒
	HeartbeatInterval int      `json:"heartbeatInterval,optional" yaml:"HeartbeatInterval"` // 单位：秒
	SessionTimeout    int      `json:"sessionTimeout,optional" yaml:"SessionTimeout"`       // 单位：秒
	RebalanceTimeout  int      `json:"rebalanceTimeout,optional" yaml:"RebalanceTimeout"`   // 单位：秒
	OffsetNewest      bool     `json:"offsetNewest,optional" yaml:"OffsetNewest"`
	Version           string   `json:"version,optional" yaml:"Version"`
	KeepAlive         int64    `json:"keepAlive,optional" yaml:"KeepAlive"` // 单位：秒

	GroupId     string         `json:"groupId,optional" yaml:"GroupId"`
	RedisConfig goredis.Config `json:"redisConfig,optional" yaml:"RedisConfig"`
}
