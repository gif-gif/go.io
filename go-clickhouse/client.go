package goclickhouse

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	gocrons "github.com/gif-gif/go.io/go-cron"
	golog "github.com/gif-gif/go.io/go-log"
)

type client struct {
	conf Config
	db   *sql.DB
	cron *gocrons.CronsModel
}

func New(conf Config) (cli *client, err error) {
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 10
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 20
	}

	cli = &client{conf: conf}

	if err = cli.connect(); err != nil {
		return
	}

	if conf.AutoPing {
		cron := gocrons.New()
		cron.Start()
		cron.SecondX(5, __client.ping)
	}

	return
}

func (cli *client) connect() (err error) {
	dns := fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s&read_timeout=%d&write_timeout=%d&alt_hosts=%s&debug=%v",
		cli.conf.Addr, cli.conf.User, cli.conf.Password, cli.conf.Database,
		cli.conf.ReadTimeout, cli.conf.WriteTimeout, cli.conf.AltHosts, cli.conf.Debug)
	cli.db, err = sql.Open(cli.conf.Driver, dns)
	if err != nil {
		golog.WithTag("goo-clickhouse").Error(err)
	}
	return
}

func (cli *client) ping() {
	if cli.db == nil {
		return
	}

	err := cli.db.Ping()
	if err == nil {
		return
	}

	if exception, ok := err.(*clickhouse.Exception); ok {
		golog.WithTag("goo-clickhouse").WithField("err_code", exception.Code).WithField("stack_trace", exception.StackTrace).Error(exception.Message)
		return
	}

	golog.WithTag("goo-clickhouse").Error(err)
}
