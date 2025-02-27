package main

import (
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	golog "github.com/gif-gif/go.io/go-log"
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

	cmd := goredis.GetClient("aa").Set("goredis", "goredis")
	if cmd.Err() != nil {
		golog.WithTag("goredis").Error(cmd.Err())
	}
	v := goredis.Default().Get("goredis").Val()
	golog.WithTag("goredis").InfoF(v)
}
