package goasynq

var __clients = map[string]*GoAsynqClient{}
var __servers = map[string]*GoAsynqServer{}

// client
func InitClient(configs ...ClientConfig) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}
		__clients[name] = NewClient(conf)
	}

	return
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

// server
func InitServer(configs ...ServerConfig) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}
		__servers[name] = RunServer(conf)
	}

	return
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
