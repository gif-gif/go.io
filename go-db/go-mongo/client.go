package gomongo

import (
	"context"
	"fmt"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//var MongoDSN = "mongodb://root:ab15eb8e12ea@122.228.113.230:27017"

type GoMongo struct {
	*mongo.Client
	conf Config
	ctx  context.Context
}

func New(conf Config) (cli *GoMongo, err error) {
	cli = &GoMongo{conf: conf, ctx: context.TODO()}

	var uri string
	if conf.EnablePassword {
		uri = fmt.Sprintf("mongodb://%s:%s@%s", conf.User, conf.Password, conf.Addr)
	} else {
		uri = fmt.Sprintf("mongodb://%s", conf.Addr)
	}

	opts := options.Client().ApplyURI(uri)
	cli.Client, err = mongo.Connect(cli.ctx, opts)
	if err != nil {
		golog.WithTag("gomongo").Error(err)
		return
	}

	if err = cli.Ping(cli.ctx, readpref.Primary()); err != nil {
		golog.WithTag("gomongo").Error(err)
		return
	}

	if conf.AutoPing {
		gj, _ := gojob.New()
		gj.Start()
		gj.SecondX(nil, 5, func() {
			if err := cli.Ping(cli.ctx, readpref.Primary()); err != nil {
				golog.WithTag("gomongo").Error(err)
			}
		})
	}

	return
}

func (cli *GoMongo) WithContext(ctx context.Context) *GoMongo {
	cli.ctx = ctx
	return cli
}

func (cli *GoMongo) DB() *mongo.Database {
	return cli.Database(cli.conf.Database)
}
