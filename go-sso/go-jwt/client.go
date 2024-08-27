package gojwt

import golog "github.com/gif-gif/go.io/go-log"

var __clients = map[string]*GoJwt{}

// 可以一次初始化多个Redis实例或者 多次调用初始化多个实例
func Init(configs ...Config) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		__clients[name] = New(conf)
	}
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
