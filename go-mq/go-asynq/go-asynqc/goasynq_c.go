package goasynqc

import (
	"errors"
	goasynq "github.com/gif-gif/go.io/go-mq/go-asynq"
)

var __clients = map[string]*goasynq.GoAsynqClient{}
var __servers = map[string]*goasynq.GoAsynqServer{}

// client
func InitClient(configs ...ClusterClientConfig) error {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("client already exists")
		}

		__clients[name] = NewClusterClient(conf)
	}

	return nil
}

func GetClient(names ...string) *goasynq.GoAsynqClient {
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

func DefaultClient() *goasynq.GoAsynqClient {
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

func DelClient(names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			delete(__clients, name)
		}
	}
}

// server
func InitServer(configs ...ClusterServerConfig) error {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}
		if __clients[name] != nil {
			return errors.New("client already exists")
		}

		__servers[name] = ClusterRunServer(conf)
	}

	return nil
}

func GetServer(names ...string) *goasynq.GoAsynqServer {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	if cli, ok := __servers[name]; ok {
		return cli
	}
	return nil
}

func DelServer(names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			delete(__servers, name)
		}
	}
}

func DefaultServer() *goasynq.GoAsynqServer {
	if cli, ok := __servers["default"]; ok {
		return cli
	}

	if l := len(__servers); l == 1 {
		for _, cli := range __servers {
			return cli
		}
	}
	return nil
}
