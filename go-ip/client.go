package goip

import (
	"errors"
)

var __clients = map[string]*GoIp{}

// 可以一次初始化多个Redis实例或者 多次调用初始化多个实例
func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("client already exists")
		}

		__clients[name], err = New(conf)
		if err != nil {
			return
		}
	}

	return
}

func GetClient(names ...string) *GoIp {
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

func Default() *GoIp {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	return nil
}
