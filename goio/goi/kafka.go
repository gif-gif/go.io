package goi

import (
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
)

func Kafka(names ...string) *gokafka.GoKafka {
	return gokafka.GetClient(names...)
}
