package gokafka

import (
	"github.com/IBM/sarama"
	gocontext "github.com/gif-gif/go.io/go-context"
	goutils "github.com/gif-gif/go.io/go-utils"
)

var (
	__client *client
)

func Init(conf Config) error {
	__client = &client{conf: conf}
	goutils.AsyncFunc(func() {
		select {
		case <-gocontext.Cancel().Done():
			__client.Close()
			return
		}
	})
	return __client.init()
}

func Client() *client {
	return __client
}

func Producer() iProducer {
	return &producer{client: __client, msg: &sarama.ProducerMessage{}}
}

func Consumer() iConsumer {
	return &consumer{client: __client}
}
