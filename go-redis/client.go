package goredis

import (
	"context"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/redis/go-redis/v9"
)

type GoRedis struct {
	*redis.Client
}

func New(conf Config) (cli *GoRedis, err error) {
	cli = &GoRedis{}

	cli.Client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
	ctx := context.Background()
	if err = cli.Ping(ctx).Err(); err != nil {
		golog.WithTag("goredis").Error(err)
		return
	}

	if conf.AutoPing {
		gj, _ := gojob.New()
		gj.Start()
		gj.SecondX(nil, 5, func() {
			if err := cli.Ping(ctx).Err(); err != nil {
				golog.WithTag("goredis").Fatal("redis ping error:", err)
			}
		})
	}

	return
}
