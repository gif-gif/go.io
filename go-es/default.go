package goes

import (
	"github.com/olivere/elastic"
)

var __client *GoEs

func Init(conf Config, options ...elastic.ClientOptionFunc) (err error) {
	__client, err = New(conf, options...)
	return
}

func Client() *GoEs {
	return __client
}

func ESClient() *elastic.Client {
	return __client.cli
}
