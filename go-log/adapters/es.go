package adapters

import (
	goes "github.com/gif-gif/go.io/go-db/go-es"
	"github.com/gif-gif/go.io/go-log"
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
	"log"
	"sync"
)

type EsAdapter struct {
	topic string
	opt   goes.Config
	mu    sync.Mutex
}

func NewEsLog(topic string, opt goes.Config) *golog.Logger {
	return golog.New(NewEsAdapter(topic, opt))
}

func NewEsAdapter(topic string, opt goes.Config) *EsAdapter {
	err := goes.Init(opt)
	if err != nil {
		return nil
	}
	fa := &EsAdapter{
		topic: topic,
		opt:   opt,
	}
	return fa
}

func (fa *EsAdapter) Write(msg *golog.Message) {
	client := gokafka.GetClient(fa.opt.Name)
	if client == nil {
		return
	}

	err := client.Producer().SendAsyncMessage(fa.topic, msg.JSON(), func(msg *gokafka.ProducerMessage, err error) {

	})

	if err != nil {
		log.Println(err.Error())
	}
}

func (fa *EsAdapter) closeEs() {
	client := gokafka.GetClient(fa.opt.Name)
	if client == nil {
		return
	}

	client.Close()
}
