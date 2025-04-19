package goasynq

import (
	"encoding/json"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type ClientConfig struct {
	Config goredis.Config `yaml:"Config" json:"config,optional"`
	Name   string         `yaml:"Name" json:"name,optional"`
}

type GoAsynqClient struct {
	Client *asynq.Client
}

func NewClient(config ClientConfig) *GoAsynqClient {
	if config.Config.PoolSize == 0 {
		config.Config.PoolSize = 10
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:         config.Config.Addr,
		DB:           config.Config.DB,
		Password:     config.Config.Password,
		PoolSize:     config.Config.PoolSize,
		DialTimeout:  lo.If(config.Config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.DialTimeout) * time.Second),
		ReadTimeout:  lo.If(config.Config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.ReadTimeout) * time.Second),
		WriteTimeout: lo.If(config.Config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.WriteTimeout) * time.Second),
	})

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
	return c.Client.Close()
}
