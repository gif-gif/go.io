package main

import (
	"encoding/json"
	golog "github.com/gif-gif/go.io/go-log"
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
	goutils "github.com/gif-gif/go.io/go-utils"
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
		Addrs:    []string{"127.0.0.1:30094"},
		User:     "admin",
		Password: "b36da6b4eb0f3",
		Timeout:  10,
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

	topic := "biu_account"
	//err = gokafka.Client().CreateTopicsRequest(topic, -1, -1)
	//if err != nil {
	//	golog.WithTag("CreateTopicsRequest").Error(err.Error())
	//	return
	//}

	goutils.AsyncFunc(func() {
		gokafka.Client().Consumer().Consume(topic, func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
			golog.WithTag("gokafka").Info("Consumer:" + msg.Topic)
			return nil
		})
	})

	_, _, err = gokafka.Client().Producer().WithPartition(0).SendMessage(topic, b)
	if err != nil {
		golog.WithTag("gokafka").Error(err.Error())
		return
	}

	golog.WithTag("gokafka").InfoF("send successfully")

	time.Sleep(time.Second * 4)
}
