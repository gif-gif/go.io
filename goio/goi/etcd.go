package goi

import (
	goetcd "github.com/gif-gif/go.io/go-etcd"
)

func Etcd(names ...string) *goetcd.GoEtcdClient {
	return goetcd.GetClient(names...)
}
