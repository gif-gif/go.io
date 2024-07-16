package main

import (
	goredis2 "github.com/gif-gif/go.io/go-db/go-redis"
	golog "github.com/gif-gif/go.io/go-log"
	"time"
)

func main() {
	config := goredis2.Config{
		Name:     "goredis",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "goredis",
		AutoPing: true,
	}

	c, err := goredis2.New(config)
	if err != nil {
		golog.WithTag("goredis").Error(err)
	}

	cmd := c.Set("goredis", "goredis", time.Duration(10)*time.Second)
	if cmd.Err() != nil {
		golog.WithTag("goredis").Error(cmd.Err())
	}
	v := c.Get("goredis").Val()
	golog.WithTag("goredis").InfoF(v)
}
