package gokafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/samber/lo"
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

func (t *KafkaMsg) Topic() string {
	return t.KafkaTopic
}

func (t *KafkaMsg) Key() string {
	return lo.If(t.KafkaTopic != "", t.KafkaTopic).Else(t.KafkaKey)
}

func (t *KafkaMsg) Headers() map[string]string {
	return lo.If(t.KafkaHeaders != nil, t.KafkaHeaders).Else(map[string]string{})
}

func (t *KafkaMsg) Serialize() []byte {
	b, _ := json.Marshal(t.Data)
	return b
}

func (t *KafkaMsg) Deserialize(b []byte) {
	if err := json.Unmarshal(b, &t.Data); err != nil {
		fmt.Println(err)
	}
}

type KafkaMsg struct {
	KafkaTopic   string            `json:"kafkaTopic"`
	KafkaKey     string            `json:"kafkaKey"`
	KafkaHeaders map[string]string `json:"kafkaHeaders"`
	Data         interface{}       `json:"data"`
}
