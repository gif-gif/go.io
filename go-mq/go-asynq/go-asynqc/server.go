package goasynqc

import (
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	golog "github.com/gif-gif/go.io/go-log"
	goasynq "github.com/gif-gif/go.io/go-mq/go-asynq"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type ClusterServerConfig struct {
	Name string
	goredisc.Config
	PoolSize    int
	Concurrency int //default 10 指定要使用的并发工作线程数量
	Queues      map[string]int
}

func ClusterRunServer(config ClusterServerConfig) *goasynq.GoAsynqServer {
	if config.Concurrency == 0 {
		config.Concurrency = 10
	}

	if config.PoolSize == 0 {
		config.PoolSize = 10
	}

	if config.Queues == nil {
		config.Queues = map[string]int{
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
			Concurrency: config.Concurrency,
			// Optionally specify multiple queues with different priority.
			Queues: config.Queues,
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	gs := &goasynq.GoAsynqServer{
		ServeMux: mux,
		Server:   srv,
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
