package goes

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
)

type GoEs struct {
	es  *elasticsearch.Client
	ctx context.Context
}

func New(options elasticsearch.Config) (cli *GoEs, err error) {
	es, err := elasticsearch.NewClient(options)

	if err != nil {
		golog.WithTag("goes").WithField("hosts", options.Addresses).WithField("options", options).Error(err)
		return nil, err
	}

	cli = &GoEs{
		ctx: context.Background(),
		es:  es,
	}

	goutils.AsyncFunc(func() {
		cli.es.Ping = func(o ...func(*esapi.PingRequest)) (*esapi.Response, error) {
			return nil, nil
		}
	})

	return cli, nil
}

func (cli *GoEs) Client() *elasticsearch.Client {
	return cli.es
}
