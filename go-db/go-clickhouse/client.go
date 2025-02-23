package goclickhouse

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/pkg/errors"
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
			//"max_execution_time": conf.MaxExecutionTime, //60,
			"max_query_size": 104857600, //100M
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
	op.Debugf = func(format string, v ...any) {
		golog.Debug(v...)
	}

	//if conf.Debugf == nil {
	//	op.Debugf = func(format string, v ...any) {
	//		golog.Debug(v...)
	//	}
	//} else {
	//	op.Debugf = conf.Debugf
	//}

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

func (cli *GoClickHouse) Close() error {
	if cli.db != nil {
		err := cli.db.Close()
		if err != nil {
			return err
		}
	}
	if cli.cron != nil {
		err := cli.cron.Stop()
		if err != nil {
			return err
		}
	}
	if cli.conn != nil {
		err := cli.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
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

// OPTIMIZE
func (cli *GoClickHouse) OptimizePartition(context context.Context, tableName string, isFinal bool, partitions ...string) error {
	finalStr := ""
	if isFinal {
		finalStr = "FINAL"
	}
	s := "OPTIMIZE TABLE " + tableName + " " + finalStr + ";"
	if len(partitions) > 0 {
		for _, partition := range partitions {
			s = "OPTIMIZE TABLE " + tableName + " PARTITION '" + partition + "' " + finalStr + ";"
			err := cli.conn.Exec(context, s)
			if err != nil {
				return err
			}
		}
	} else {
		err := cli.conn.Exec(context, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// 注意：
// 物化视图时需要表名特殊处理如下
//
// mvInnerTableName := ".inner_id.2f31ec0d-f667-487e-8b18-c13324c27bc7" 用于获取分组列表
//
// tableName := "`" + mvInnerTableName + "`" 用于处理具体的分区
func (cli *GoClickHouse) GetPartitions(ctx context.Context, tableName string) ([]PartitionInfo, error) {
	query := `
		SELECT 
			partition,
			sum(bytes) as size
		FROM system.parts 
		WHERE table = $1 
		GROUP BY partition 
		ORDER BY partition
	`

	rows, err := cli.conn.Query(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("查询分区信息失败: %v", err)
	}
	defer rows.Close()

	var partitions []PartitionInfo
	for rows.Next() {
		var p PartitionInfo
		if err := rows.Scan(&p.Partition, &p.Size); err != nil {
			return nil, fmt.Errorf("读取分区信息失败: %v", err)
		}
		partitions = append(partitions, p)
	}

	return partitions, nil
}

// 获取分区内的最小和最大ID
func (cli *GoClickHouse) GetPartitionIDRange(ctx context.Context, tableName, partition string) (min, max uint64, err error) {
	query := fmt.Sprintf(`
		SELECT 
			min(id) as min_id,
			max(id) as max_id
		FROM %s 
		WHERE partition = $1
	`, tableName)

	row := cli.conn.QueryRow(ctx, query, partition)
	if err := row.Scan(&min, &max); err != nil {
		return 0, 0, fmt.Errorf("获取ID范围失败: %v", err)
	}
	return min, max, nil
}

// 删除分区内的部分数据
func (cli *GoClickHouse) DeletePartitionData(ctx context.Context, tableName string, partition string, startID, endID uint64) error {
	query := fmt.Sprintf(`
		ALTER TABLE %s 
		DELETE WHERE partition = $1 AND id >= $2 AND id < $3
	`, tableName)

	if err := cli.conn.Exec(ctx, query, partition, startID, endID); err != nil {
		return fmt.Errorf("删除数据失败: %v", err)
	}
	return nil
}

// 删除整个分区(如果分区小于50GB，直接删除整个分区 时可以直接删除,否则用 DropPartition)
func (cli *GoClickHouse) dropPartition(ctx context.Context, tableName, partition string) error {
	query := fmt.Sprintf(`ALTER TABLE %s DROP PARTITION $1`, tableName)
	if err := cli.conn.Exec(ctx, query, partition); err != nil {
		return fmt.Errorf("删除分区失败: %v", err)
	}
	return nil
}

// 删除分区(支持大分区删除)
//
// 注意：
// 物化视图时需要表名特殊处理如下
//
// mvInnerTableName := ".inner_id.2f31ec0d-f667-487e-8b18-c13324c27bc7" 用于获取分组列表
//
// tableName := "`" + mvInnerTableName + "`" 用于处理具体的分区
func (cli *GoClickHouse) DropPartition(ctx context.Context, tableName string, p PartitionInfo) error {
	if p.Size <= MaxSizeToDelete {
		// 如果分区小于50GB，直接删除整个分区
		if err := cli.dropPartition(ctx, tableName, p.Partition); err != nil {
			return errors.Wrapf(err, "删除分区 %s 失败: %v", p.Partition)
		}
		return nil
	} else {
		// 如果分区大于50GB，分批删除数据
		minID, maxID, err := cli.GetPartitionIDRange(ctx, tableName, p.Partition)
		if err != nil {
			return errors.Wrapf(err, "获取分区 %s 的ID范围失败: %v", p.Partition)
		}

		for startID := minID; startID < maxID; startID += BatchSizeToDelete {
			endID := startID + BatchSizeToDelete
			if endID > maxID {
				endID = maxID
			}

			if err := cli.DeletePartitionData(ctx, tableName, p.Partition, startID, endID); err != nil {
				return errors.Wrapf(err, "删除分区 %s 中的数据(ID范围: %d-%d)失败: %v", p.Partition, startID, endID)
			}

			//log.Printf("成功删除分区 %s 中的数据(ID范围: %d-%d)", p.Partition, startID, endID)
			// 可选：添加一些延时避免过度占用系统资源
			time.Sleep(time.Second * 1)
			return nil
		}
	}

	return nil
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
