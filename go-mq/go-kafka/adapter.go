package gokafka

import "github.com/IBM/sarama"

// 生产者
type iProducer interface {
	init() error

	Close()

	Client() sarama.Client

	// 发送消息到指定分区
	WithPartition(partition int32) iProducer

	// 发送消息 - 同步
	SendMessage(topic string, message []byte) (partition int32, offset int64, err error)

	// 发送消息 - 异步
	SendAsyncMessage(topic string, message []byte, cb MessageHandler) (err error)
}

// 消费者
type iConsumer interface {
	init() error

	Close()

	Client() sarama.Client

	// 从指定分区消费
	WithPartition(partition int32) iConsumer

	// 从指定位置开始
	WithOffset(offset int64) iConsumer

	// 从最新位置开始
	WithOffsetNewest() iConsumer

	// 从头开始
	WithOffsetOldest() iConsumer

	// 消费
	Consume(topic string, handler ConsumerHandler)

	// 分组topic
	ConsumeGroup(groupId string, topics []string, handler ConsumerHandler)
}
