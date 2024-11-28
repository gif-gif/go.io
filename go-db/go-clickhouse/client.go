package goclickhouse

import (
	"context"
	"crypto/tls"
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"time"
)

type GoClickHouse struct {
	conf Config
	db   *sql.DB
	conn driver.Conn
	cron *gojob.GoJob
}

func New(conf Config) (cli *GoClickHouse, err error) {
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

	cli = &GoClickHouse{conf: conf}

	getTls := &tls.Config{
		InsecureSkipVerify: conf.InsecureSkipVerify,
	}

	if !conf.TLS {
		getTls = nil
	}
	op := &clickhouse.Options{
		Addr: conf.Addr,
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
		//Settings: clickhouse.Settings{
		//	"max_memory_usage":                 "10000000000",  // 增加内存使用限制
		//	"max_bytes_before_external_group_by": "20000000000", // 增加分组操作前的字节限制
		//	"max_block_size":                   "100000",       // 调整块大小
		//},
		DialTimeout: time.Second * time.Duration(conf.DialTimeout),
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                conf.Debug,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see GoClickHouse info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "go.io", Version: "0.1"},
			},
		},
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
	}

	if conf.Debugf == nil {
		op.Debugf = func(format string, v ...any) {
			golog.Debug(v...)
		}
	} else {
		op.Debugf = conf.Debugf
	}

	cli.db = clickhouse.OpenDB(op)
	conn, err := clickhouse.Open(op)
	if err != nil {
		return nil, err
	}
	cli.conn = conn

	cli.db.SetMaxIdleConns(conf.MaxIdleConn)
	cli.db.SetMaxOpenConns(conf.MaxOpenConn)
	cli.db.SetConnMaxLifetime(time.Second * time.Duration(conf.ConnMaxLifetime))

	if conf.AutoPing {
		cron, err := gojob.New()
		if err != nil {
			return nil, err
		}
		cron.Start()
		_, err = cron.SecondX(nil, 5, cli.ping)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (cli *GoClickHouse) DB() *sql.DB {
	return cli.db
}

func (cli *GoClickHouse) Conn() driver.Conn {
	return cli.conn
}

func (cli *GoClickHouse) ping() {
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

// BaseModel 注意T必须为指针类型
type BaseModel[T any] struct {
	Client driver.Conn
	Table  string
}

// BatchInsert 注意添加字段时，先发布代码，再往数据库添加字段。不然先加字段会出现插不进去
func (m *BaseModel[T]) BatchInsert(ctx context.Context, items []T) error {
	batch, err := m.Client.PrepareBatch(ctx, "INSERT INTO "+m.Table)
	if err != nil {
		return err
	}
	for i := range items {
		err := batch.AppendStruct(items[i])
		if err != nil {
			return err
		}
	}
	err = batch.Send()
	return err
}
