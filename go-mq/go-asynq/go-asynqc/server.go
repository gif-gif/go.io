package goasynqc

import (
	"context"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"time"
)

type GoAsynqServer struct {
	ServeMux *asynq.ServeMux
	Server   *asynq.Server
}

type ClusterServerConfig struct {
	Name        string          `yaml:"Name" json:"name,optional"`
	Config      goredisc.Config `yaml:"Config" json:"config,optional"`
	Concurrency int             `yaml:"Concurrency" json:"concurrency,optional"` //default 10 指定要使用的并发工作线程数量
	Queues      map[string]int  `yaml:"Queues" json:"queues,optional"`
}

func ClusterRunServer(conf ClusterServerConfig) *GoAsynqServer {
	config := conf.Config
	server := &asynq.Server{}
	if config.Type != "cluster" {
		server = RunServerByNode(conf)
	} else {
		server = RunServer(conf)
	}

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	gs := &GoAsynqServer{
		ServeMux: mux,
		Server:   server,
	}

	goutils.AsyncFunc(func() { // 异步运行挂起
		defer gs.Stop()
		if err := server.Run(mux); err != nil {
			golog.WithTag("GoAsynqServer").Error("could not run server: %v", err)
		}
		golog.WithTag("GoAsynqServer").Info("stop running")
	})

	return gs
}

func RunServer(conf ClusterServerConfig) *asynq.Server {
	if conf.Concurrency == 0 {
		conf.Concurrency = 10
	}

	if conf.Queues == nil {
		conf.Queues = map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		}
	}
	config := conf.Config
	server := asynq.NewServer(
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

	return server
}

func RunServerByNode(config ClusterServerConfig) *asynq.Server {
	if config.Config.PoolSize == 0 {
		config.Config.PoolSize = 10
	}

	if config.Concurrency == 0 {
		config.Concurrency = 10
	}

	if config.Queues == nil {
		config.Queues = map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		}
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:         config.Config.Addrs[0],
			Password:     config.Config.Password,
			DB:           config.Config.DB,
			PoolSize:     config.Config.PoolSize,
			DialTimeout:  lo.If(config.Config.DialTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.DialTimeout) * time.Second),
			ReadTimeout:  lo.If(config.Config.ReadTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.ReadTimeout) * time.Second),
			WriteTimeout: lo.If(config.Config.WriteTimeout <= 0, time.Duration(5)*time.Second).Else(time.Duration(config.Config.WriteTimeout) * time.Second),
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: config.Concurrency,
			// Optionally specify multiple queues with different priority.
			Queues: config.Queues,
			// See the godoc for other configuration options
		},
	)

	return srv
}

func (s *GoAsynqServer) HandleFunc(taskTypeTopic string, handler func(context.Context, *asynq.Task) error) {
	s.ServeMux.HandleFunc(taskTypeTopic, handler)
}

func (s *GoAsynqServer) Handle(taskTypeTopic string, handler asynq.Handler) {
	s.ServeMux.Handle(taskTypeTopic, handler)
}

// Stop指示服务器停止从队列中提取新任务。
// 在关闭服务器之前，可以使用Stop来确保所有
// 在服务器关闭之前处理当前活动的任务。
//
// Stop不会关闭服务器，请确保在退出前调用shutdown。
func (s *GoAsynqServer) Stop() {
	s.Server.Stop()
}

// Shutdown会优雅地关闭服务器。
// 它优雅地关闭了所有活跃的员工。服务器将等待
// 在配置中指定的持续时间内，主动工作人员完成处理任务。关机超时。
// 如果worker在超时期间没有完成任务处理，则该任务将被推回Redis。
func (s *GoAsynqServer) Shutdown() {
	s.Server.Shutdown()
}
