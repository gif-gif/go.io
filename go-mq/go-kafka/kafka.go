package gokafka

import (
	"errors"
	"github.com/IBM/sarama"
	gocontext "github.com/gif-gif/go.io/go-context"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
)

var __clients = map[string]*client{}

func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" || name == "default" {
			conf.Name = "default"
		}
		if __clients[name] != nil {
			return errors.New("client already exists")
		}
		__clients[name], err = New(conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func New(conf Config) (*client, error) {
	__client := &client{conf: conf}
	goutils.AsyncFunc(func() {
		select {
		case <-gocontext.Cancel().Done():
			__client.Close()
			return
		}
	})
	err := __client.init()
	return __client, err
}

func GetClient(names ...string) *client {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	if cli, ok := __clients[name]; ok {
		return cli
	}
	return nil
}

func DelClient(names ...string) {
	if l := len(names); l > 0 {
		for _, name := range names {
			delete(__clients, name)
		}
	}
}

// default or 只有一个kafka实例直接返回
func Client() *client {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("gokafka").Error("no default kafka client")

	return nil
}

func Producer() iProducer {
	return &producer{client: Client(), msg: &sarama.ProducerMessage{}}
}

func Consumer() iConsumer {
	return &consumer{client: Client()}
}
