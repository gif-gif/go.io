package goasynqc

import (
	"errors"
)

var __clients = map[string]*GoAsynqClient{}
var __servers = map[string]*GoAsynqServer{}
var __inspector = map[string]*GoAsynqInspector{}

// client for cluster or node
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

func GetClient(names ...string) *GoAsynqClient {
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

func DefaultClient() *GoAsynqClient {
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

// server for cluster or node
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

func GetServer(names ...string) *GoAsynqServer {
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

func DefaultServer() *GoAsynqServer {
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

// __inspector for cluster or node
func InitInspector(configs ...ClusterInspectorConfig) error {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __inspector[name] != nil {
			return errors.New("__inspector already exists")
		}

		__inspector[name] = NewClusterInspector(conf)
	}

	return nil
}

func GetInspector(names ...string) *GoAsynqInspector {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	if cli, ok := __inspector[name]; ok {
		return cli
	}
	return nil
}

func DelInspector(names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			delete(__inspector, name)
		}
	}
}

func DefaultInspector() *GoAsynqInspector {
	if cli, ok := __inspector["default"]; ok {
		return cli
	}

	if l := len(__inspector); l == 1 {
		for _, cli := range __inspector {
			return cli
		}
	}
	return nil
}
