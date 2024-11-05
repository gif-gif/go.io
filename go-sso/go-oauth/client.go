package gooauth

import (
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
)

var __clients = map[string]*GoOAuth{}

// 可以一次初始化多个实例或者 多次调用初始化多个实例
func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("GoAuth [" + name + "] already exists")
		}

		__clients[name] = New(conf)
		if err != nil {
			return
		}
	}

	return
}

func GetClient(names ...string) *GoOAuth {
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
			delete(__clients, name)
		}
	}
}

func Default() *GoOAuth {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("GoAuth").Error("no default GoAuth client")

	return nil
}