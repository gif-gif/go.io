package gokafka

import (
	"github.com/IBM/sarama"
	gocontext "github.com/gif-gif/go.io/go-context"
)

type MessageHandler func(msg *ProducerMessage, err error)

type ProducerMessage struct {
	*sarama.ProducerMessage
}

type ConsumerHandler func(ctx *gocontext.Context, msg *ConsumerMessage, consumerErr *ConsumerError) error

type ConsumerMessage struct {
	*sarama.ConsumerMessage
	GroupSession sarama.ConsumerGroupSession
}

func (msg ConsumerMessage) Commit() {
	if msg.GroupSession == nil {
		return
	}
	msg.GroupSession.MarkMessage(msg.ConsumerMessage, "")
}

type ConsumerError struct {
	*sarama.ConsumerError
}
