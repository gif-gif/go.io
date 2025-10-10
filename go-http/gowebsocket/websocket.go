package gowebsocket

import (
	"errors"
	"sync"
	"time"

	golog "github.com/gif-gif/go.io/go-log"
)

var mu sync.RWMutex
var __servers = map[string]*Server{}
var __clients = map[string]*WSClient{}

// server--------------

// 可以一次初始化多个Redis实例或者 多次调用初始化多个实例
func InitServer(configs ...ServerConfig) (err error) {
	for _, conf := range configs {
		mu.Lock()
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __servers[name] != nil {
			return errors.New("Hub server [" + name + "] already exists")
		}

		__servers[name] = NewServer(conf)
		mu.Unlock()
	}

	return nil
}

func GetServer(names ...string) *Server {
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
			r := GetServer(name)
			if r != nil {
				if r.hub != nil {
					err := r.GracefulShutdown(30 * time.Second)
					if err != nil {
						golog.WithTag("GracefulShutdown").Error("err:" + err.Error())
					}
				}
			}
			delete(__servers, name)
		}
	}
}

func DefaultServer() *Server {
	if cli, ok := __servers["default"]; ok {
		return cli
	}

	if l := len(__servers); l == 1 {
		for _, cli := range __servers {
			return cli
		}
	}
	golog.WithTag("Hub").Error("no default websocket server")
	return nil
}

// client ------------

// 可以一次初始化多个Redis实例或者 多次调用初始化多个实例
func InitClient(configs ...ClientConfig) (err error) {
	for _, conf := range configs {
		mu.Lock()
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("Hub client [" + name + "] already exists")
		}

		__clients[name] = NewWSClient(conf)
		mu.Unlock()
	}

	return nil
}

func GetClient(names ...string) *WSClient {
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
			r := GetClient(name)
			if r != nil {
				r.Close()
			}
			delete(__clients, name)
		}
	}
}

func DefaultClient() *WSClient {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}
	golog.WithTag("Hub").Error("no default websocket client")
	return nil
}
