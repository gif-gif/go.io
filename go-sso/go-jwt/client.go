package gojwt

import (
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
)

var __clients = map[string]*GoJwt{}

func Init(configs ...Config) error {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("gojwt client [" + name + "] already exists")
		}

		__clients[name] = New(conf)
	}

	return nil
}

func GetClient(names ...string) *GoJwt {
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

func Default() *GoJwt {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("gojwt").Error("no default jwt client")
	return nil
}
