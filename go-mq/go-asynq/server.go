package goasynq

import (
	"context"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/hibiken/asynq"
	"log"
)

type ServerConfig struct {
	Name string
	goredis.Config
	PoolSize    int
	Concurrency int //default 10 指定要使用的并发工作线程数量
	Queues      map[string]int
}

type GoAsynqServer struct {
	ServeMux *asynq.ServeMux
	Server   *asynq.Server
}

func (s *GoAsynqServer) Stop() {
	s.Server.Stop()
}

func RunServer(config ServerConfig) *GoAsynqServer {
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
		asynq.RedisClientOpt{Addr: config.Addr, Password: config.Password, DB: config.DB, PoolSize: config.PoolSize},
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

	gs := &GoAsynqServer{
		ServeMux: mux,
		Server:   srv,
	}

	goutils.AsyncFunc(func() { // 异步运行挂起
		defer gs.Stop()
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
		golog.WithTag("GoAsynqServer").Info("stop running")
	})

	return gs
}

func (s *GoAsynqServer) HandleFunc(taskTypeTopic string, handler func(context.Context, *asynq.Task) error) {
	s.ServeMux.HandleFunc(taskTypeTopic, handler)
}

func (s *GoAsynqServer) Handle(taskTypeTopic string, handler asynq.Handler) {
	s.ServeMux.Handle(taskTypeTopic, handler)
}
