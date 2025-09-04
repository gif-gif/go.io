package gogorm

import (
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/gorm"
)

var _clients = map[string]*GoGorm{}

func GetNames() []string {
	names := make([]string, 0, len(_clients))
	for name := range _clients {
		names = append(names, name)
	}
	return names
}

func InitWithGormConfig(gormConfig gorm.Config, configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if _clients[name] != nil {
			return errors.New("gogorm client already exists")
		}

		_clients[name], err = New(&conf, gormConfig)
		if err != nil {
			return
		}
	}

	return
}

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

func DelClient(names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			delete(_clients, name)
		}
	}
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
