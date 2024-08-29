package main

import (
	"encoding/json"
	gocontext "github.com/gif-gif/go.io/go-context"
	golog "github.com/gif-gif/go.io/go-log"
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
	goutils "github.com/gif-gif/go.io/go-utils"
	"gorm.io/gorm"
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
		Addrs:        []string{"122.228.113.231:30094"}, //122.228.113.231
		User:         "admin",
		Password:     "b36da6b4eb0f3",
		Timeout:      10,
		OffsetNewest: false,
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
	err = gokafka.Client().CreateTopicsRequest(topic, 0, -1)
	if err != nil {
		golog.WithTag("CreateTopicsRequest").Error(err.Error())
		return
	}

	//goutils.AsyncFunc(func() {
	//	gokafka.Consumer().Consume(topic, func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
	//		golog.WithTag("gokafka").Info("Consumer:" + msg.Topic)
	//		golog.WithTag("gokafka").Info("Consumer:", string(msg.Value))
	//		return nil
	//	})
	//})
	//
	goutils.AsyncFunc(func() {
		gokafka.Consumer().ConsumeGroup("pro", []string{topic}, func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
			golog.WithTag("gokafkaGroup").Info("Consumer:" + msg.Topic)
			golog.WithTag("gokafkaGroup").Info("Consumer:", string(msg.Value))
			return nil
		})
	})

	_, _, err = gokafka.Producer().WithPartition(0).SendMessage(topic, b)
	if err != nil {
		golog.WithTag("gokafka").Error(err.Error())
		return
	}

	golog.WithTag("gokafka").InfoF("send successfully")

	<-gocontext.Cancel().Done()
}
