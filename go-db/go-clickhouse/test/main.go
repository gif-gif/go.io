package main

import (
	"fmt"
	goclickhouse2 "github.com/gif-gif/go.io/go-db/go-clickhouse"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"log"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	err := goclickhouse2.Init(goclickhouse2.Config{
		Driver:             "clickhouse",
		Addr:               "122.228.113.238:8124",
		User:               "default",
		Password:           "payda6b4eb0f3",
		Database:           "xzdsp",
		InsecureSkipVerify: true,
	})

	if err != nil {
		golog.Fatal(err)
		return
	}
	rows, err := goclickhouse2.DB().Query("SELECT oaid FROM xzdsp.clickcb limit 10")
	if err != nil {
		golog.Fatal(err)
		return
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
