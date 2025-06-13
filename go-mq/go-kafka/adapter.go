package gokafka

import "github.com/IBM/sarama"

// 生产者
type IProducer interface {
	Client() sarama.Client

	// 发送消息 - 同步
	SendMessage(msg IMessage) (partition int32, offset int64, err error)

	SendMessages(msgs []IMessage) (err error)

	// 发送消息 - 异步
	SendAsyncMessage(msg IMessage, cb MessageHandler) (err error)
	SendAsyncMessages(msgs []IMessage, cb MessageHandler) (err error)
}

// 消费者
type IConsumer interface {
	Client() sarama.Client

	// 从指定分区消费
	WithPartition(partition int32) IConsumer

	// 从指定位置开始
	WithOffset(offset int64) IConsumer

	// 从最新位置开始
	WithOffsetNewest() IConsumer

	// 从头开始
	WithOffsetOldest() IConsumer

	// 消费
	Consume(topic string, handler ConsumerHandler)

	// 分组topic
	ConsumeGroup(groupId string, topics []string, handler ConsumerHandler)
}

// 消息
type IMessage interface {
	Topic() string
	Key() string
	Headers() map[string]string
	Serialize() []byte
	Deserialize(b []byte)
}
