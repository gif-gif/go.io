package goasynqc

import (
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type GoAsynqInspector struct {
	Inspector *asynq.Inspector
}

type ClusterInspectorConfig struct {
	Config goredisc.Config `yaml:"Config" json:"config,optional"`
	Name   string          `yaml:"Name" json:"name,optional"`
}

func NewClusterInspector(conf ClusterInspectorConfig) *GoAsynqInspector {
	config := conf.Config
	var inspector *asynq.Inspector
	if config.Type != "cluster" {
		inspector = asynq.NewInspector(asynq.RedisClientOpt{
			Addr:         config.Addrs[0],
			Password:     config.Password,
			DB:           config.DB,
			PoolSize:     config.PoolSize,
			DialTimeout:  lo.If(config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.DialTimeout) * time.Second),
			ReadTimeout:  lo.If(config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.ReadTimeout) * time.Second),
			WriteTimeout: lo.If(config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.WriteTimeout) * time.Second),
		})
	} else {
		inspector = asynq.NewInspector(asynq.RedisClusterClientOpt{
			Addrs:        config.Addrs,
			Password:     config.Password,
			DialTimeout:  lo.If(config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.DialTimeout) * time.Second),
			ReadTimeout:  lo.If(config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.ReadTimeout) * time.Second),
			WriteTimeout: lo.If(config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.WriteTimeout) * time.Second),
		})
	}

	return &GoAsynqInspector{
		Inspector: inspector,
	}
}

func (c *GoAsynqInspector) Close() error {
	err := c.Inspector.Close()
	if err != nil {
		return err
	}

	//
	//if c.Redis != nil {
	//	return c.Redis.Close()
	//}

	return nil
}

//func (c *GoAsynqInspector) DeleteQueueFromRedis(queueName string) error {
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
func (c *GoAsynqInspector) PurgeQueue(queueName string) error {
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

	return c.DeleteQueue(queueName, true)
}

func (c *GoAsynqInspector) DeleteQueue(queueName string, force bool) error {
	return c.Inspector.DeleteQueue(queueName, force)
}

func (c *GoAsynqInspector) Queues() ([]string, error) {
	infos, err := c.Inspector.Queues()
	if err != nil {
		return nil, err
	}
	return infos, nil
}
