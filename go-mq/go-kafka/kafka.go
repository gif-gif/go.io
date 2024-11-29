package gokafka

import (
	"errors"
	gocontext "github.com/gif-gif/go.io/go-context"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
)

var __clients = map[string]*GoKafka{}

func Init(conf Config, opts ...Option) (err error) {
	if conf.Name == "" {
		conf.Name = "default"
	}
	if __clients[conf.Name] != nil {
		return errors.New("GoKafka already exists")
	}
	__clients[conf.Name], err = New(conf, opts...)
	if err != nil {
		return err
	}

	return nil
}

func New(conf Config, opts ...Option) (*GoKafka, error) {
	__client := &GoKafka{conf: conf}
	goutils.AsyncFunc(func() {
		select {
		case <-gocontext.WithCancel().Done():
			__client.Close()
			return
		}
	})
	for _, opt := range opts {
		switch opt.Name {
		case RedisName:
			__client.redis = opt.Value.(*goredis.GoRedis)
		}
	}

	err := __client.init()
	return __client, err
}

func GetClient(names ...string) *GoKafka {
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
func Client() *GoKafka {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("gokafka").Error("no default kafka GoKafka")

	return nil
}

func Consumer() IConsumer {
	return GetClient().Consumer()
}

func Producer(opts ...Option) IProducer {
	return GetClient().Producer(opts...)
}
