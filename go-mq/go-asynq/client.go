package goasynq

import (
	"encoding/json"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	"github.com/hibiken/asynq"
)

type ClientConfig struct {
	goredis.Config
	PoolSize int
}

type GoAsynqClient struct {
	Client *asynq.Client
}

func NewClient(config ClientConfig) *GoAsynqClient {

	if config.PoolSize == 0 {
		config.PoolSize = 10
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     config.Addr,
		DB:       config.DB,
		Password: config.Password,
		PoolSize: config.PoolSize,
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
