# kafka
- https://github.com/IBM/sarama
```
gokafka.Client().Topics()

// 发送消息，不指定分区
gokafka.Producer().SendMessage("test", []byte("hi"))

// 发送消息，指定分区
gokafka.Producer().WithPartition(0).SendMessage("test", []byte("hi goio"))

// 发送异步消息，不指定分区
gokafka.Producer().SendAsyncMessage("test", []byte("hi goio"), func(msg *gokafka.ProducerMessage, err error) {
})

// 发送异步消息，指定分区
gokafka.Producer().WithPartition(0).SendAsyncMessage("test", []byte("hi goio"), func(msg *gokafka.ProducerMessage, err error) {
})

// 消费消息，指定分区，指定起始位置
gokafka.Consumer().WithPartition(0).WithOffset(100).Consume("test", func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
    return nil
})

// 消费消息，指定分区，从最新位置开始
gokafka.Consumer().WithPartition(0).WithOffsetNewest().Consume("test", func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
    return nil
})

// 消费消息，指定分区，从最头开始
gokafka.Consumer().WithPartition(0).WithOffsetOldest().Consume("test", func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
    return nil
})

// 消费消息，分组消息，分组里面只要1个消费者消费
gokafka.Consumer().ConsumeGroup("test-id", []string{"test"}, func(msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
    return nil
})
```
