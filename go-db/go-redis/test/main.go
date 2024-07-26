package main

import (
	"context"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	golog "github.com/gif-gif/go.io/go-log"
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

	err := goredis.Init(config)
	if err != nil {
		golog.WithTag("goredis").Error(err)
	}

	cmd := goredis.Default().Set(context.Background(), "goredis", "goredis", time.Duration(10)*time.Second)
	if cmd.Err() != nil {
		golog.WithTag("goredis").Error(cmd.Err())
	}
	v := goredis.Default().Get(context.Background(), "goredis").Val()
	golog.WithTag("goredis").InfoF(v)
}
