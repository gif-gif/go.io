package adapters

import (
	"github.com/gif-gif/go.io/go-log"
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
	"sync"
)

type KafkaAdapter struct {
	topic string
	opt   gokafka.Config
	mu    sync.Mutex
}

func NewKafkaLog(topic string, opt gokafka.Config) *golog.Logger {
	return golog.New(NewKafkaAdapter(topic, opt))
}

func NewKafkaAdapter(topic string, opt gokafka.Config) *KafkaAdapter {
	err := gokafka.Init(opt)
	if err != nil {
		return nil
	}
	fa := &KafkaAdapter{
		topic: topic,
		opt:   opt,
	}
	return fa
}

func (fa *KafkaAdapter) Write(msg *golog.Message) {
	client := gokafka.GetClient(fa.opt.Name)
	if client == nil {
		return
	}

	//err := client.Producer().SendAsyncMessage(msg.JSON(), func(msg *gokafka.ProducerMessage, err error) {
	//
	//})
	//
	//if err != nil {
	//	log.Println(err.Error())
	//}
}

func (fa *KafkaAdapter) CloseKafka() {
	client := gokafka.GetClient(fa.opt.Name)
	if client == nil {
		return
	}

	client.Close()
}
