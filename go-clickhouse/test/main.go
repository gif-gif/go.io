package main

import (
	"fmt"
	goclickhouse "github.com/gif-gif/go.io/go-clickhouse"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"log"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	goclickhouse.Init(goclickhouse.Config{
		Driver:             "clickhouse",
		Addr:               "122.228.113.238:8124",
		User:               "default",
		Password:           "payda6b4eb0f3",
		Database:           "xzdsp",
		InsecureSkipVerify: true,
	})

	rows, err := goclickhouse.DB().Query("SELECT oaid FROM xzdsp.clickcb limit 10")
	if err != nil {
		golog.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			oaid string
		)
		if err := rows.Scan(&oaid); err != nil {
			log.Fatal(err)
		}
		fmt.Println(oaid)
	}

	if err := rows.Err(); err != nil {
		golog.Fatal(err)
	}
}
