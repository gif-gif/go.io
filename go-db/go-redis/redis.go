package goredis

import (
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
)

var __clients = map[string]*GoRedis{}

// 可以一次初始化多个Redis实例或者 多次调用初始化多个实例
func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("goredis client [" + name + "] already exists")
		}

		__clients[name], err = New(conf)
		if err != nil {
			return err
		}
	}

	return
}

func GetClient(names ...string) *GoRedis {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	if cli, ok := __clients[name]; ok {
		return cli
	}
	return nil
}

func DelClient(names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			r := GetClient(name)
			if r != nil {
				r.Close()
			}

			delete(__clients, name)
		}
	}
}

func Default() *GoRedis {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("goredis").Error("no default redis client")

	return nil
}
