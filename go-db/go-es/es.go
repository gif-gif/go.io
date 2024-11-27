package goes

import (
	"errors"
	"github.com/olivere/elastic/v7"
)

var __clients = map[string]*GoEs{}

// 可以多次调用初始化多个实例
func Init(conf Config, options ...elastic.ClientOptionFunc) (err error) {
	name := conf.Name
	if name == "" {
		name = "default"
	}

	if __clients[name] != nil {
		return errors.New("GoEs already exists")
	}

	__clients[name], err = New(conf, options...)
	if err != nil {
		return
	}

	return nil
}

func GetClient(names ...string) *GoEs {
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

func Default() *GoEs {
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
