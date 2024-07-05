package gomongo

import golog "github.com/gif-gif/go.io/go-log"

var __clients = map[string]*Client{}

func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		__clients[name], err = New(conf)
		if err != nil {
			return
		}
	}

	return
}

func GetClient(names ...string) *Client {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
	}

	if cli, ok := __clients[name]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("gomongo").Error("no default mongo client")

	return nil
}

func Default() *Client {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("gomongo").Error("no default mongo client")

	return nil
}
