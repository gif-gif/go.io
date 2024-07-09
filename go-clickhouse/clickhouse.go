package goclickhouse

import "database/sql"

var __client *GoClickHouse

func DB() *sql.DB {
	return __client.db
}

// 全局
func Init(conf Config) error {
	client, err := CreateConnection(conf)
	if err != nil {
		return err
	}

	__client = client
	return nil
}

// 创建一个连接
func CreateConnection(conf Config) (*GoClickHouse, error) {
	client, err := New(conf)
	if err != nil {
		return nil, err
	}

	return client, nil
}
