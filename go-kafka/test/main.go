package main

import (
	"encoding/json"
	gokafka "github.com/gif-gif/go.io/go-kafka"
	golog "github.com/gif-gif/go.io/go-log"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}
type Account struct {
	Id        int64 `json:"id"`
	UpdatedAt int64 `json:"updated_at"`
}

func main() {
	err := gokafka.Init(gokafka.Config{
		Addrs:    []string{"212.129.60.103:30092"},
		User:     "admin",
		Password: "payda6b4eb0f3",
	})

	if err != nil {
		golog.WithTag("gokafka").Error(err.Error())
		return
	}
	
	msg := Account{}
	b, err := json.Marshal(msg)
	if err != nil {
		golog.WithTag("gokafka").Error(err.Error())
		return
	}

	_, _, err = gokafka.Producer().SendMessage("biu_account", b)
	if err != nil {
		golog.WithTag("gokafka").Error(err.Error())
		return
	}

	golog.WithTag("gokafka").InfoF("send successfully")

	time.Sleep(time.Second * 10000)
}
