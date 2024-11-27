package goi

import (
	goes "github.com/gif-gif/go.io/go-db/go-es"
	"github.com/olivere/elastic/v7"
)

func EsClient(names ...string) *goes.GoEs {
	return goes.GetClient(names...)
}

func ES(names ...string) *elastic.Client {
	client := EsClient(names...)
	return client.Client()
}
