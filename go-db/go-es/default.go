package goes

import "github.com/olivere/elastic/v7"

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
