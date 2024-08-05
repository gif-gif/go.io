package main

import (
	"encoding/json"
	golog "github.com/gif-gif/go.io/go-log"
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
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

	gokafka.GetClient().Consumer().Consume("biu_account", func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
		golog.WithTag("gokafka").Info(msg.Topic)
		return nil
	})

	_, _, err = gokafka.GetClient().Producer().SendMessage("biu_account", b)
	if err != nil {
		golog.WithTag("gokafka").Error(err.Error())
		return
	}

	golog.WithTag("gokafka").InfoF("send successfully")

	time.Sleep(time.Second * 10000)
}
