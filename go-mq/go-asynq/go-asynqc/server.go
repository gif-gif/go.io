package goasynqc

import (
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	golog "github.com/gif-gif/go.io/go-log"
	goasynq "github.com/gif-gif/go.io/go-mq/go-asynq"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type ClusterServerConfig struct {
	Name        string          `yaml:"Name" json:"name,optional"`
	Config      goredisc.Config `yaml:"Config" json:"config,optional"`
	Concurrency int             `yaml:"Concurrency" json:"concurrency,optional"` //default 10 指定要使用的并发工作线程数量
	Queues      map[string]int  `yaml:"Queues" json:"queues,optional"`
	Prefix      string          `yaml:"Prefix" json:"prefix,optional"`
}

func convertServerConfigToNode(conf *ClusterServerConfig) goasynq.ServerConfig {
	return goasynq.ServerConfig{
		Config: goredis.Config{
			Name:         conf.Config.Name,
			Addr:         conf.Config.Addrs[0],
			DB:           conf.Config.DB,
			Password:     conf.Config.Password,
			Prefix:       conf.Config.Prefix,
			TLS:          conf.Config.TLS,
			AutoPing:     conf.Config.AutoPing,
			PoolSize:     conf.Config.PoolSize,
			DialTimeout:  conf.Config.DialTimeout,
			ReadTimeout:  conf.Config.ReadTimeout,
			WriteTimeout: conf.Config.WriteTimeout,
			Type:         "node",
			Weight:       conf.Config.Weight,
		},
		Queues: conf.Queues,
	}
}

func ClusterRunServer(conf ClusterServerConfig) *goasynq.GoAsynqServer {
	config := conf.Config
	if config.Type != "cluster" {
		return goasynq.RunServer(convertServerConfigToNode(&conf))
	}
	if conf.Concurrency == 0 {
		conf.Concurrency = 10
	}

	if config.PoolSize == 0 {
		config.PoolSize = 10
	}

	if conf.Queues == nil {
		conf.Queues = map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		}
	}

	srv := asynq.NewServer(
		asynq.RedisClusterClientOpt{
			Addrs:        config.Addrs,
			Password:     config.Password,
			DialTimeout:  lo.If(config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.DialTimeout) * time.Second),
			ReadTimeout:  lo.If(config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.ReadTimeout) * time.Second),
			WriteTimeout: lo.If(config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.WriteTimeout) * time.Second),
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: conf.Concurrency,
			// Optionally specify multiple queues with different priority.
			Queues: conf.Queues,
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	gs := &goasynq.GoAsynqServer{
		ServeMux: mux,
		Server:   srv,
		Prefix:   config.Prefix,
	}

	goutils.AsyncFunc(func() { // 异步运行挂起
		defer gs.Stop()
		if err := srv.Run(mux); err != nil {
			golog.WithTag("GoAsynqServer").Error("could not run server: %v", err)
		}
		golog.WithTag("GoAsynqServer").Info("stop running")
	})

	return gs
}
