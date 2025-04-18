package goasynqc

import (
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	goasynq "github.com/gif-gif/go.io/go-mq/go-asynq"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type ClusterClientConfig struct {
	goredisc.Config
	Name string `yaml:"Name" json:"name,optional"`
}

func NewClusterClient(config ClusterClientConfig) *goasynq.GoAsynqClient {
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
