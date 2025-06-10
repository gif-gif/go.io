package goasynq

import (
	"encoding/json"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type ClientConfig struct {
	Config goredisc.Config `yaml:"Config" json:"config,optional"`
	Name   string          `yaml:"Name" json:"name,optional"`
}

type GoAsynqClient struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
	//Redis     *goredisc.GoRedisC
}

func NewClient(config ClientConfig) *GoAsynqClient {
	if config.Config.PoolSize == 0 {
		config.Config.PoolSize = 10
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:         config.Config.Addrs[0],
		DB:           config.Config.DB,
		Password:     config.Config.Password,
		PoolSize:     config.Config.PoolSize,
		DialTimeout:  lo.If(config.Config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.DialTimeout) * time.Second),
		ReadTimeout:  lo.If(config.Config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.ReadTimeout) * time.Second),
		WriteTimeout: lo.If(config.Config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.WriteTimeout) * time.Second),
	})

	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr:         config.Config.Addrs[0],
		Password:     config.Config.Password,
		DB:           config.Config.DB,
		PoolSize:     config.Config.PoolSize,
		DialTimeout:  lo.If(config.Config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.DialTimeout) * time.Second),
		ReadTimeout:  lo.If(config.Config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.ReadTimeout) * time.Second),
		WriteTimeout: lo.If(config.Config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.WriteTimeout) * time.Second),
	})

	//err := goredisc.Init(config.Config)
	//if err != nil {
	//	logx.Errorf("init redis client error: %s", err)
	//}

	return &GoAsynqClient{
		Client:    client,
		Inspector: inspector,
		//Redis:     goredisc.GetClient(config.Name),
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

func (c *GoAsynqClient) Queues() ([]string, error) {
	infos, err := c.Inspector.Queues()
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (c *GoAsynqClient) DeleteQueue(queueName string, force bool) error {
	return c.Inspector.DeleteQueue(queueName, force)
}

func (c *GoAsynqClient) Close() error {
	err := c.Inspector.Close()
	if err != nil {
		return err
	}

	err = c.Client.Close()

	if err != nil {
		return err
	}
	//
	//if c.Redis != nil {
	//	return c.Redis.Close()
	//}

	return nil
}

//func (c *GoAsynqClient) DeleteQueueFromRedis(queueName string) error {
//	// 创建 Redis 客户端
//	rdb := c.Redis
//	// Asynq 使用的 Redis key 模式
//	patterns := []string{
//		fmt.Sprintf("asynq:{%s}:*", queueName),
//		fmt.Sprintf("asynq:queues:%s", queueName),
//		// 其他相关的 key 模式
//	}
//
//	for _, pattern := range patterns {
//		keys, err := rdb.Keys(pattern).Result()
//		if err != nil {
//			return err
//		}
//
//		if len(keys) > 0 {
//			err = rdb.Del(keys...).Err()
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}

// 彻底地清理队列
func (c *GoAsynqClient) PurgeQueue(queueName string) error {
	inspector := c.Inspector

	// 获取队列信息
	queues, err := inspector.Queues()
	if err != nil {
		return err
	}

	// 检查队列是否存在
	queueExists := false
	for _, q := range queues {
		if q == queueName {
			queueExists = true
			break
		}
	}

	if !queueExists {
		return nil
	}

	// 删除各种状态的任务
	deleteFuncs := []func(string) (int, error){
		inspector.DeleteAllPendingTasks,
		inspector.DeleteAllRetryTasks,
		inspector.DeleteAllScheduledTasks,
		inspector.DeleteAllArchivedTasks,
		inspector.DeleteAllCompletedTasks,
	}

	for _, deleteFunc := range deleteFuncs {
		if _, err := deleteFunc(queueName); err != nil {
			// 继续删除其他状态的任务
			return err
		}
	}

	return nil
}
