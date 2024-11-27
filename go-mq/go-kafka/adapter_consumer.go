package gokafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	gocontext "github.com/gif-gif/go.io/go-context"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
)

type consumer struct {
	*GoKafka

	hasSetPartition bool  // 是否设置分区
	partition       int32 // 分区

	offset int64
}

func (c *consumer) Client() sarama.Client {
	return c.GoKafka.Client
}

// 设置 分区
func (c *consumer) WithPartition(partition int32) iConsumer {
	c.hasSetPartition = true
	c.partition = partition
	return c
}

// 设置 起始位置
func (c *consumer) WithOffset(offset int64) iConsumer {
	c.offset = offset
	return c
}

// 设置 起始位置 = 最新位置
func (c *consumer) WithOffsetNewest() iConsumer {
	c.offset = sarama.OffsetNewest
	return c
}

// 设置 起始位置 = 从头开始
func (c *consumer) WithOffsetOldest() iConsumer {
	c.offset = sarama.OffsetOldest
	return c
}

// 消费消息，默认处理最新消息
func (c *consumer) Consume(topic string, handler ConsumerHandler) {
	consumer, err := sarama.NewConsumerFromClient(c.Client())
	if err != nil {
		golog.WithTag("gokafka-consumer").Error(err)
		return
	}
	defer consumer.Close()

	if c.offset == 0 {
		c.offset = sarama.OffsetNewest
	}

	pc, err := consumer.ConsumePartition(topic, c.partition, c.offset)
	if err != nil {
		golog.WithTag("gokafka-consumer").Error(err)
		return
	}
	defer pc.Close()

	for {
		select {
		case <-gocontext.WithCancel().Done():
			return

		case msg := <-pc.Messages():
			handler(&ConsumerMessage{msg}, nil)

		case err := <-pc.Errors():
			handler(nil, &ConsumerError{err})
		}
	}
}

// 分组
func (c *consumer) ConsumeGroup(groupId string, topics []string, handler ConsumerHandler) {
	cg, err := sarama.NewConsumerGroupFromClient(groupId, c.Client())
	if err != nil {
		golog.WithTag("gokafka-consumer-group").Error(err)
		return
	}

	g := group{handler: handler}
	goutils.AsyncFunc(func() {
		defer cg.Close()
		for {
			select {
			case <-gocontext.WithCancel().Done():
				return

			case err := <-cg.Errors():
				golog.WithTag("gokafka-consumer-group").Error(err)
				return
			}
		}
	})

	if err := cg.Consume(gocontext.WithCancel().Context, topics, g); err != nil {
		golog.WithTag("gokafka-consumer-group").Error(err)
	}
}

func (c *consumer) ConsumeGroup1(groupId string, topics []string, handler ConsumerHandler) {
	g, err := sarama.NewConsumerGroup(c.conf.Addrs, groupId, c.Client().Config())
	if err != nil {
		panic(err)
	}
	defer func() { _ = g.Close() }()

	// Track errors
	go func() {
		for err := range g.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		ghandler := group{handler: handler}
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := g.Consume(ctx, topics, ghandler)
		if err != nil {
			panic(err)
		}
	}
}
