package goredis

import (
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/go-redis/redis"
)

type Client struct {
	*redis.Client
}

func New(conf Config) (cli *Client, err error) {
	cli = &Client{}

	cli.Client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	if err = cli.Ping().Err(); err != nil {
		golog.WithTag("goredis").Error(err)
		return
	}

	if conf.AutoPing {
		gj, _ := gojob.New()
		gj.Start()
		gj.SecondX(nil, 5, func() {
			if err := cli.Ping().Err(); err != nil {
				golog.WithTag("goredis").Error(err)
			}
		})
	}

	return
}
