package goi

import (
	goip "github.com/gif-gif/go.io/go-ip"
)

func GoIp(names ...string) *goip.GoIp {
	return goip.GetClient(names...)
}
