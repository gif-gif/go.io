package goi

import (
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
)

func Kafka(names ...string) *gokafka.GoKafka {
	return gokafka.GetClient(names...)
}

func Producer(names ...string) gokafka.IProducer {
	return gokafka.GetClient(names...).Producer()
}

func Consumer(names ...string) gokafka.IConsumer {
	return gokafka.GetClient(names...).Consumer()
}
