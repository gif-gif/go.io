package goclickhouse

import (
	"database/sql"
	"time"
)

var __client *client

func Init(conf Config) {
	var err error
	if __client, err = New(conf); err != nil {
		time.Sleep(10 * time.Second)
		Init(conf)
	}
}

func DB() *sql.DB {
	return __client.db
}
