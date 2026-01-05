package gokafka

import (
	"context"
	"errors"
	"time"

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
func (c *consumer) WithPartition(partition int32) IConsumer {
	c.hasSetPartition = true
	c.partition = partition
	return c
}

// 设置 起始位置
func (c *consumer) WithOffset(offset int64) IConsumer {
	c.offset = offset
	return c
}

// 设置 起始位置 = 最新位置
func (c *consumer) WithOffsetNewest() IConsumer {
	c.offset = sarama.OffsetNewest
	return c
}

// 设置 起始位置 = 从头开始
func (c *consumer) WithOffsetOldest() IConsumer {
	c.offset = sarama.OffsetOldest
	return c
}

// 消费消息，默认处理最新消息
func (c *consumer) Consume(topic string, handler ConsumerHandler) {
	log := golog.WithTag("gokafka-consumer").WithField("topic", topic)

	consumer, err := sarama.NewConsumerFromClient(c.Client())
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error(err)
		}
	}()

	if c.offset == 0 {
		c.offset = sarama.OffsetNewest
	}

	pc, err := consumer.ConsumePartition(topic, c.partition, c.offset)
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		if err := pc.Close(); err != nil {
			log.Error(err)
		}
	}()

	ctx := gocontext.WithCancel()

	for {
		select {
		case <-ctx.Done():
			log.Debug("Context被取消,停止消费")
			return

		case err := <-pc.Errors():
			if err != nil {
				log.Error(err)
			}

		case msg, ok := <-pc.Messages():
			if !ok {
				log.Debug("消息通道被关闭,停止消费")
				return
			}

			onceCtx, cancel := context.WithCancel(ctx.Context)
			handlerCtx := gocontext.WithParent(onceCtx).WithLog()
			handlerCtx.Log.WithTag("gokafka-consumer").WithField("msg", msg)

			if err = handler(handlerCtx, &ConsumerMessage{ConsumerMessage: msg}, nil); err != nil {
				log.Error(err)
			}

			// 删除缓存
			key := string(msg.Key)
			if c.redis != nil {
				c.redis.Del(key)
			}

			cancel()
		}
	}
}

// 分组
func (c *consumer) ConsumeGroup(groupId string, topics []string, handler ConsumerHandler) {
	l := golog.WithTag("gokafka-consumer-group").
		WithField("groupId", groupId).
		WithField("topics", topics)

	cg, err := sarama.NewConsumerGroupFromClient(groupId, c.GoKafka.Client)
	if err != nil {
		l.Error(err)
		return
	}
	defer func() {
		if c.GoKafka.Client != nil {
			c.GoKafka.Client.Close()
			l.Debug("client 退出")
		}
	}()
	defer func() {
		cg.Close()
		l.Debug("consumer-group 退出")
	}()

	var (
		done = make(chan struct{})
		flag bool
	)

	ctx, cancel := context.WithCancel(context.Background())

	goutils.AsyncFunc(func() {
		for {
			select {
			case err := <-cg.Errors():
				if err != nil {
					l.Error(err)
				}

			default:
				err := cg.Consume(ctx, topics, group{id: groupId, handler: handler, GoKafka: c.GoKafka})
				if err != nil && !errors.Is(err, sarama.ErrClosedConsumerGroup) {
					l.Error(err)
				}
				if flag {
					done <- struct{}{}
					return
				}
			}
		}
	})

	select {
	case <-gocontext.WithCancel().Done():
		flag = true
		cancel()
	}

	<-done

	time.Sleep(time.Second)
}
