package gogorm

import (
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
)

var _clients = map[string]*GoGorm{}

func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if _clients[name] != nil {
			return errors.New("gogorm client already exists")
		}

		_clients[name], err = New(&conf)
		if err != nil {
			return
		}
	}

	return
}

func GetClient(names ...string) *GoGorm {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	if cli, ok := _clients[name]; ok {
		return cli
	}
	return nil
}

func Default() *GoGorm {
	if cli, ok := _clients["default"]; ok {
		return cli
	}

	if l := len(_clients); l == 1 {
		for _, cli := range _clients {
			return cli
		}
	}

	golog.WithTag("gogorm").Error("no default gogorm client")

	return nil
}
