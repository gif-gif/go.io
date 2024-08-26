package gomqtt

import (
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
)

var __clients = map[string]*GoMqttClient{}

func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" || name == "default" {
			conf.Name = "default"
		}

		if GetClient(conf.Name) != nil {
			return errors.New(conf.Name + " 实例已存在")
		}
		__clients[name], err = NewClient(conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetClient(names ...string) *GoMqttClient {
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

func Client() *GoMqttClient {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("gokafka").Error("no default kafka client")

	return nil
}
