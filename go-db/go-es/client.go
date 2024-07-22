package goes

import (
	"context"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/olivere/elastic/v7"
)

type client struct {
	cli *elastic.Client
	ctx context.Context
}

func New(conf Config, options ...elastic.ClientOptionFunc) (cli *client, err error) {
	options = append(options,
		// 将sniff设置为false后，便不会自动转换地址
		elastic.SetSniff(false),
		elastic.SetGzip(true),
		elastic.SetURL(conf.Addr),
		elastic.SetBasicAuth(conf.User, conf.Password),
	)

	if conf.EnableLog {
		options = append(options, elastic.SetErrorLog(logger{golog.ERROR}))
		options = append(options, elastic.SetTraceLog(logger{golog.DEBUG}))
		options = append(options, elastic.SetInfoLog(logger{golog.INFO}))
	}

	cli = &client{
		ctx: context.Background(),
	}

	cli.cli, err = elastic.NewClient(options...)
	if err != nil {
		golog.WithTag("goo-es").WithField("host", conf.Addr).WithField("options", options).Error(err)
		return
	}

	goutils.AsyncFunc(func() {
		_, _, err := cli.cli.Ping(conf.Addr).Do(cli.ctx)
		if err != nil {
			golog.WithTag("goo-es").WithField("host", conf.Addr).WithField("options", options).Error(err)
			return
		}
	})

	goutils.AsyncFunc(func() {
		_, err := cli.cli.ElasticsearchVersion(conf.Addr)
		if err != nil {
			golog.WithTag("goo-es").WithField("host", conf.Addr).WithField("options", options).Error(err)
			return
		}
	})

	return
}

func (cli *client) Client() *elastic.Client {
	return cli.cli
}
