package main

import (
	golog "github.com/gif-gif/go.io/go-log"
	goredis "github.com/gif-gif/go.io/go-redis"
	"time"
)

func main() {
	config := goredis.Config{
		Name:     "goredis",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "goredis",
		AutoPing: true,
	}

	c, err := goredis.New(config)
	if err != nil {
		golog.WithTag("goredis").Error(err)
	}

	cmd := c.Set("goredis", "goredis", time.Duration(0))
	if cmd.Err() != nil {
		golog.WithTag("goredis").Error(cmd.Err())
	}
	v := c.Get("goredis").Val()
	golog.WithTag("goredis").InfoF(v)
}
