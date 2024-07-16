package gokafka

import "github.com/IBM/sarama"

type MessageHandler func(msg *ProducerMessage, err error)

type ProducerMessage struct {
	*sarama.ProducerMessage
}

type ConsumerHandler func(msg *ConsumerMessage, consumerErr *ConsumerError) error

type ConsumerMessage struct {
	*sarama.ConsumerMessage
}

type ConsumerError struct {
	*sarama.ConsumerError
}
