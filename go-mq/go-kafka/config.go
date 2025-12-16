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

	AutoCommit     AutoCommitOption `json:"autoCommit" yaml:"AutoCommit"`
	ProducerFlush  ProducerFlush    `json:"producerFlush" yaml:"ProducerFlush"`
	ConsumerConfig ConsumerConfig   `json:"consumerConfig" yaml:"ConsumerConfig"`

	ChannelBufferSize int `json:"ChannelBufferSize,default=1024" yaml:"ChannelBufferSize,default=1024"`
}

type AutoCommitOption struct {
	Enable   bool `json:"enable,default=true" yaml:"Enable,default=true"`
	Interval int  `json:"interval,default=1" yaml:"Interval,default=1"` // 单位：秒
}

type ProducerFlush struct {
	// 生产者累计消息数发送，默认100
	Messages int `json:"messages,default=100" yaml:"Messages,default=100"`
	// 生产者累计消息大小，默认1MB
	Bytes int `json:"bytes,default=100" yaml:"Bytes,default=100"`
	// 生产者每x秒提交，默认100ms
	Frequency int `json:"frequency,default=100" yaml:"Frequency,default=100"`
}

type (
	ConsumerConfig struct {
		GroupConfig         ConsumerGroupConfig `json:"groupConfig" yaml:"GroupConfig"`
		ConsumerFetchConfig ConsumerFetchConfig `json:"consumerFetchConfig" yaml:"ConsumerFetchConfig"`
	}

	ConsumerGroupConfig struct {
		// 心跳监测，默认5秒
		HeartbeatInterval int `json:"heartbeatInterval,default=5" yaml:"HeartbeatInterval,default=5"`
		// session超时时间 默认15秒
		SessionTimeout int `json:"sessionTimeout,default=15" yaml:"SessionTimeout,default=15"`
		// Reblance 超时时间 moren 12秒
		ReblanceInterval int `json:"rebalanceInterval,default=12" yaml:"RebalanceInterval,default=12"`
	}

	ConsumerFetchConfig struct {
		// 默认抓取数量 15MB
		Default int32 `json:"default,default=15" yaml:"Default,default=15"`
		// 最小抓取数量 100 Kb
		Min int32 `json:"min,default=100" yaml:"Min,default=100"`
		// 最大抓取 50MB
		Max int32 `json:"max,default=15" yaml:"Max,default=15"`
	}
)
