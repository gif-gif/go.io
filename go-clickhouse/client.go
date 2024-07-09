package goclickhouse

import (
	"crypto/tls"
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"time"
)

type Client struct {
	conf Config
	db   *sql.DB
	cron *gojob.GoJob
}

func New(conf Config) (cli *Client, err error) {

	if conf.DialTimeout == 0 {
		conf.DialTimeout = 60
	}

	if conf.MaxIdleConn == 0 {
		conf.MaxIdleConn = 5
	}
	if conf.MaxOpenConn == 0 {
		conf.MaxOpenConn = 10
	}
	if conf.ConnMaxLifetime == 0 {
		conf.ConnMaxLifetime = 60 * 60
	}

	cli = &Client{conf: conf}

	getTls := &tls.Config{
		InsecureSkipVerify: conf.InsecureSkipVerify,
	}

	if !conf.Tls {
		getTls = nil
	}

	cli.db = clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{conf.Addr},
		Auth: clickhouse.Auth{
			Database: conf.Database,
			Username: conf.User,
			Password: conf.Password,
		},
		Protocol: clickhouse.HTTP,
		TLS:      getTls,
		Settings: clickhouse.Settings{
			"max_execution_time": conf.MaxExecutionTime, //60,
		},
		DialTimeout: time.Second * time.Duration(conf.DialTimeout),
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug: conf.Debug,
		Debugf: func(format string, v ...any) {
			golog.Debug(v...)
		},
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "go.io", Version: "0.1"},
			},
		},
	})

	cli.db.SetMaxIdleConns(conf.MaxIdleConn)
	cli.db.SetMaxOpenConns(conf.MaxOpenConn)
	cli.db.SetConnMaxLifetime(time.Second * time.Duration(conf.ConnMaxLifetime))

	if conf.AutoPing {
		cron, err := gojob.New()
		if err != nil {
			return nil, err
		}
		cron.Start()
		_, err = cron.SecondX(nil, 5, __client.ping)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (cli *Client) ping() {
	if cli.db == nil {
		return
	}

	err := cli.db.Ping()
	if err == nil {
		return
	}

	if exception, ok := err.(*clickhouse.Exception); ok {
		golog.WithTag("goclickhouse").WithField("err_code", exception.Code).WithField("stack_trace", exception.StackTrace).Error(exception.Message)
		return
	}

	golog.WithTag("goclickhouse").Error(err)
}
