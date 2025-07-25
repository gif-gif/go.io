package goasynqc

import (
	"encoding/json"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type GoAsynqClient struct {
	Client *asynq.Client
}

type ClusterClientConfig struct {
	Config goredisc.Config `yaml:"Config" json:"config,optional"`
	Name   string          `yaml:"Name" json:"name,optional"`
}

func NewClusterClient(conf ClusterClientConfig) *GoAsynqClient {
	config := conf.Config
	var client *asynq.Client
	if config.Type != "cluster" {
		if conf.Config.PoolSize == 0 {
			conf.Config.PoolSize = 10
		}

		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr:         conf.Config.Addrs[0],
			DB:           conf.Config.DB,
			Password:     conf.Config.Password,
			PoolSize:     conf.Config.PoolSize,
			DialTimeout:  lo.If(conf.Config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(conf.Config.DialTimeout) * time.Second),
			ReadTimeout:  lo.If(conf.Config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(conf.Config.ReadTimeout) * time.Second),
			WriteTimeout: lo.If(conf.Config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(conf.Config.WriteTimeout) * time.Second),
		})
	} else {
		client = asynq.NewClient(asynq.RedisClusterClientOpt{
			Addrs:        config.Addrs,
			Password:     config.Password,
			MaxRedirects: lo.If(config.MaxRedirects <= 0, 10).Else(config.MaxRedirects),
			DialTimeout:  lo.If(config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.DialTimeout) * time.Second),
			ReadTimeout:  lo.If(config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.ReadTimeout) * time.Second),
			WriteTimeout: lo.If(config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.WriteTimeout) * time.Second),
		})
	}

	return &GoAsynqClient{
		Client: client,
	}
}

func (c *GoAsynqClient) NewTask(taskTypeTopic string, payload any, opts ...asynq.Option) (*asynq.Task, error) {
	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(taskTypeTopic, payloadByte, opts...), nil
}

func (c *GoAsynqClient) Enqueue(taskTypeTopic string, payload any, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	task, err := c.NewTask(taskTypeTopic, payload, opts...)
	if err != nil {
		return nil, err
	}
	info, err := c.Client.Enqueue(task)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (c *GoAsynqClient) Close() error {
	err := c.Client.Close()
	if err != nil {
		return err
	}
	//
	//if c.Redis != nil {
	//	return c.Redis.Close()
	//}

	return nil
}
