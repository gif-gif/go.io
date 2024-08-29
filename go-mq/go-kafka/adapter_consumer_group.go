package gokafka

import (
	"github.com/IBM/sarama"
)

// 分组
type group struct {
	handler ConsumerHandler
}

func (g group) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (group) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g group) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := g.handler(&ConsumerMessage{msg}, nil); err != nil {
			continue
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
