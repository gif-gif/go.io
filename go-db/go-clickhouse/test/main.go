package main

import (
	"fmt"
	goclickhouse2 "github.com/gif-gif/go.io/go-db/go-clickhouse"
	golog "github.com/gif-gif/go.io/go-log"
	"log"
)

func main() {
	err := goclickhouse2.Init(goclickhouse2.Config{
		Addr:               []string{"122.28.113.238:111"},
		User:               "default",
		Password:           "111",
		Database:           "xzdsp",
		InsecureSkipVerify: true,
	})

	if err != nil {
		golog.Fatal(err)
		return
	}

	rows, err := goclickhouse2.Default().DB().Query("SELECT oaid FROM xzdsp.clickcb limit 10")
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
