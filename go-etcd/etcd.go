package goetcd

import (
	"errors"

	golog "github.com/gif-gif/go.io/go-log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var __clients = map[string]*GoEtcdClient{}

func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		if __clients[name] != nil {
			return errors.New("client already exists")
		}

		__clients[name], err = New(conf)
		if err != nil {
			return
		}
	}

	return
}

func GetClient(names ...string) *GoEtcdClient {
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

func Default() *GoEtcdClient {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("goetcd").Error("no default goetcd client")

	return nil
}

func Set(key, val string) (resp *clientv3.PutResponse, err error) {
	return Default().Set(key, val)
}

func SetWithPrevKV(key, val string) (resp *clientv3.PutResponse, err error) {
	return Default().SetWithPrevKV(key, val)
}

// ttl is seconds
func SetTTL(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	return Default().SetTTL(key, val, ttl)
}

// ttl is seconds
func SetTTLWithPrevKV(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	return Default().SetTTLWithPrevKV(key, val, ttl)
}

func Get(key string, opts ...clientv3.OpOption) (rsp *clientv3.GetResponse, err error) {
	return Default().Get(key, opts...)
}

func Exists(key string, opts ...clientv3.OpOption) (bool, error) {
	e, err := Get(key, opts...)
	if err != nil {
		return false, err
	}
	return e.Count > 0, err
}

func GetString(key string) string {
	return Default().GetString(key)
}

func GetArray(key string) (data []string) {
	return Default().GetArray(key)
}

func GetMap(key string) (data map[string]string) {
	return Default().GetMap(key)
}

func Del(key string) (resp *clientv3.DeleteResponse, err error) {
	return Default().Del(key)
}

func DelWithPrefix(key string) (resp *clientv3.DeleteResponse, err error) {
	return Default().DelWithPrefix(key)
}

func RegisterService(key, val string) (err error) {
	return Default().RegisterService(key, val)
}

func Watch(key string) <-chan []string {
	return Default().Watch(key)
}
