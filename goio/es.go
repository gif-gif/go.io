package goio

import goes "github.com/gif-gif/go.io/go-db/go-es"

func ES(names ...string) *goes.GoEs {
	return goes.GetClient(names...)
}
