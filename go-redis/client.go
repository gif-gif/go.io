package goredis

import (
	gocrons "github.com/gif-gif/go.io/go-cron"
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
		golog.WithTag("goo-redis").Error(err)
		return
	}

	if conf.AutoPing {
		gocrons.SecondX(5, func() {
			if err := cli.Ping().Err(); err != nil {
				golog.WithTag("goo-redis").Error(err)
			}
		})
	}

	return
}
