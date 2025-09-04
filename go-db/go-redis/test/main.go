package main

import (
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	golog "github.com/gif-gif/go.io/go-log"
)

func main() {
	config := goredisc.Config{
		Name:     "goredis",
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "",
		DB:       0,
		Prefix:   "goredis",
		AutoPing: true,
	}

	err := goredisc.Init(config)
	if err != nil {
		golog.WithTag("goredis").Error(err)
	}

	cmd := goredisc.Default().Set("goredis", "goredis")
	if cmd.Err() != nil {
		golog.WithTag("goredis").Error(cmd.Err())
	}
	v := goredisc.Default().Get("goredis").Val()
	golog.WithTag("goredis").InfoF(v)
}
