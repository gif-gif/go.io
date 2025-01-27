package goes

import (
	"context"
	"github.com/gif-gif/go.io/go-db/go-es/eslog"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/olivere/elastic/v7"
)

type GoEs struct {
	cli *elastic.Client
	ctx context.Context
}

func New(conf Config, options ...elastic.ClientOptionFunc) (cli *GoEs, err error) {
	options = append(options,
		// 将sniff设置为false后，便不会自动转换地址
		elastic.SetSniff(false),
		elastic.SetGzip(true),
		elastic.SetURL(conf.Addr),
		elastic.SetBasicAuth(conf.User, conf.Password),
	)

	if conf.EnableLog {
		options = append(options, elastic.SetErrorLog(eslog.Logger{golog.ERROR}))
		options = append(options, elastic.SetTraceLog(eslog.Logger{golog.DEBUG}))
		options = append(options, elastic.SetInfoLog(eslog.Logger{golog.INFO}))
	}

	cli = &GoEs{
		ctx: context.Background(),
	}

	cli.cli, err = elastic.NewClient(options...)
	if err != nil {
		golog.WithTag("go-es").WithField("host", conf.Addr).WithField("options", options).Error(err)
		return
	}

	goutils.AsyncFunc(func() {
		_, _, err := cli.cli.Ping(conf.Addr).Do(cli.ctx)
		if err != nil {
			golog.WithTag("go-es").WithField("host", conf.Addr).WithField("options", options).Error(err)
			return
		}
	})

	goutils.AsyncFunc(func() {
		_, err := cli.cli.ElasticsearchVersion(conf.Addr)
		if err != nil {
			golog.WithTag("go-es").WithField("host", conf.Addr).WithField("options", options).Error(err)
			return
		}
	})

	return
}

func (cli *GoEs) Client() *elastic.Client {
	return cli.cli
}
