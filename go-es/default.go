package goo_es

import (
	"github.com/olivere/elastic"
)

var __client *client

func Init(conf Config, options ...elastic.ClientOptionFunc) (err error) {
	__client, err = New(conf, options...)
	return
}

func Client() *client {
	return __client
}

func ESClient() *elastic.Client {
	return __client.cli
}
