package goasynqc

import (
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	goasynq "github.com/gif-gif/go.io/go-mq/go-asynq"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type ClusterClientConfig struct {
	Config goredisc.Config `yaml:"Config" json:"config,optional"`
	Name   string          `yaml:"Name" json:"name,optional"`
}

func convertClientConfigToNode(conf *ClusterClientConfig) goasynq.ClientConfig {
	config := conf.Config
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	return goasynq.ClientConfig{
		Config: goredis.Config{
			Name:         config.Name,
			Addr:         config.Addrs[0],
			DB:           config.DB,
			Password:     config.Password,
			Prefix:       config.Prefix,
			TLS:          config.TLS,
			AutoPing:     config.AutoPing,
			PoolSize:     config.PoolSize,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			Type:         "node",
			Weight:       config.Weight,
		},
		Name: config.Name,
	}
}

func NewClusterClient(conf ClusterClientConfig) *goasynq.GoAsynqClient {
	config := conf.Config
	if config.Type != "cluster" {
		return goasynq.NewClient(convertClientConfigToNode(&conf))
	}
	client := asynq.NewClient(asynq.RedisClusterClientOpt{
		Addrs:        config.Addrs,
		Password:     config.Password,
		DialTimeout:  lo.If(config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.DialTimeout) * time.Second),
		ReadTimeout:  lo.If(config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.ReadTimeout) * time.Second),
		WriteTimeout: lo.If(config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.WriteTimeout) * time.Second),
	})

	return &goasynq.GoAsynqClient{
		Client: client,
	}
}
