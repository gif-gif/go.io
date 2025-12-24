package gokafka

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"time"
)

type producerFast struct {
	*GoKafka
	focus         bool // 是否强制发送
	asyncProducer sarama.AsyncProducer
}

func (p *producerFast) Client() sarama.Client {
	return p.GoKafka.Client
}

// 发送消息 - 同步- 批量
func (p *producerFast) SendMessages(msgs []IMessage) (err error) {
	if len(msgs) == 0 {
		return nil
	}
	var sendMsgs []*sarama.ProducerMessage
	for _, msg := range msgs {
		m := &sarama.ProducerMessage{
			Topic: msg.Topic(),
			Value: sarama.ByteEncoder(msg.Serialize()),
		}
		if v := msg.Key(); v != "" {
			m.Key = sarama.StringEncoder(v)
		}

		if data := msg.Headers(); data != nil {
			var headers []sarama.RecordHeader
			for k, v := range data {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(k),
					Value: []byte(v),
				})
			}
			m.Headers = headers
		}
		sendMsgs = append(sendMsgs, m)
	}

	defer func() {
		log := golog.WithTag("gokafka-producer").WithField("msg", goutils.M{
			"send_msgs_length": len(sendMsgs),
		})
		if err != nil {
			log.Error("消息发送失败", err)
			return
		}
		//log.Debug("消息发送成功")
	}()

	var producer sarama.SyncProducer
	producer, err = sarama.NewSyncProducerFromClient(p.Client())
	if err != nil {
		return
	}
	defer producer.Close()
	err = producer.SendMessages(sendMsgs)
	return
}

// 发送消息 - 同步
func (p *producerFast) SendMessage(msg IMessage) (partition int32, offset int64, err error) {
	m := &sarama.ProducerMessage{
		Topic: msg.Topic(),
		Value: sarama.ByteEncoder(msg.Serialize()),
	}

	if v := msg.Key(); v != "" {
		m.Key = sarama.StringEncoder(v)
	}
	if data := msg.Headers(); data != nil {
		var headers []sarama.RecordHeader
		for k, v := range data {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
		m.Headers = headers
	}

	defer func() {
		log := golog.WithTag("gokafka-producer").WithField("msg", goutils.M{
			"topic":     msg.Topic(),
			"key":       msg.Key(),
			"headers":   msg.Headers(),
			"body":      msg,
			"partition": m.Partition,
			"offset":    m.Offset,
		})
		if err != nil {
			log.Error("消息发送失败", err)
			return
		}
		//log.Debug("消息发送成功")
	}()

	// 添加缓存
	if p.redis != nil && len(msg.Key()) > 0 {
		if p.focus {
			p.redis.Del(msg.Key())
		}
		if p.redis.Exists(msg.Key()).Val() > 0 {
			err = errors.New("KEY已存在")
			return
		}
		p.redis.Set1(msg.Key(), goutils.M{
			"topic":     msg.Topic(),
			"body":      msg,
			"headers":   msg.Headers(),
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		}.String(), time.Hour)
	}

	var producer sarama.SyncProducer

	producer, err = sarama.NewSyncProducerFromClient(p.Client())
	if err != nil {
		return
	}
	defer producer.Close()

	partition, offset, err = producer.SendMessage(m)

	return
}

// 在初始化该结构体的地方调用
func (p *producerFast) initAsyncEngine() {
	var err error
	p.asyncProducer, err = sarama.NewAsyncProducerFromClient(p.GoKafka.Client)
	if err != nil {
		panic(fmt.Sprintf("初始化Kafka异步引擎失败: %v", err))
	}

	// 【核心】启动后台协程处理所有消息的回执
	go func() {
		for {
			select {
			case _msg := <-p.asyncProducer.Successes():
				// 从 Metadata 中取出当时传入的回调函数
				if cb, ok := _msg.Metadata.(MessageHandler); ok && cb != nil {
					cb(&ProducerMessage{_msg}, nil)
				}
			case e := <-p.asyncProducer.Errors():
				if cb, ok := e.Msg.Metadata.(MessageHandler); ok && cb != nil {
					cb(&ProducerMessage{e.Msg}, e.Err)
				}
			}
		}
	}()
}

// 发送消息 - （全局单生产者）异步 - 批量
func (p *producerFast) SendAsyncMessages(msgs []IMessage, cb MessageHandler) (err error) {
	if len(msgs) == 0 {
		return nil
	}

	for _, msg := range msgs {
		m := &sarama.ProducerMessage{
			Topic:    msg.Topic(),
			Value:    sarama.ByteEncoder(msg.Serialize()),
			Metadata: cb,
		}

		if v := msg.Key(); v != "" {
			m.Key = sarama.StringEncoder(v)
		}

		if data := msg.Headers(); data != nil {
			var headers []sarama.RecordHeader
			for k, v := range data {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(k),
					Value: []byte(v),
				})
			}
			m.Headers = headers
		}

		p.asyncProducer.Input() <- m
	}

	return
}

// 发送消息 - 异步
func (p *producerFast) SendAsyncMessage(msg IMessage, cb MessageHandler) (err error) {
	m := &sarama.ProducerMessage{
		Topic: msg.Topic(),
		Value: sarama.ByteEncoder(msg.Serialize()),
	}

	if v := msg.Key(); v != "" {
		m.Key = sarama.StringEncoder(v)
	}
	if data := msg.Headers(); data != nil {
		var headers []sarama.RecordHeader
		for k, v := range data {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
		m.Headers = headers
	}

	defer func() {
		log := golog.WithTag("gokafka-producer").WithField("msg", goutils.M{
			"topic":     msg.Topic(),
			"key":       msg.Key(),
			"headers":   msg.Headers(),
			"body":      msg,
			"partition": m.Partition,
			"offset":    m.Offset,
		})
		if err != nil {
			log.Error("消息发送失败", err)
			return
		}
		//log.Debug("消息发送成功")
	}()

	// 添加缓存
	if p.redis != nil && len(msg.Key()) > 0 {
		if p.focus {
			p.redis.Del(msg.Key())
		}
		if p.redis.Exists(msg.Key()).Val() > 0 {
			err = errors.New("KEY已存在")
			return
		}
		p.redis.Set1(msg.Key(), goutils.M{
			"topic":     msg.Topic(),
			"body":      msg,
			"headers":   msg.Headers(),
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		}.String(), time.Hour)
	}

	var producer sarama.AsyncProducer

	producer, err = sarama.NewAsyncProducerFromClient(p.Client())
	if err != nil {
		return
	}
	defer producer.Close()

	producer.Input() <- m
	select {
	case msg := <-producer.Successes():
		cb(&ProducerMessage{msg}, nil)
	case e := <-producer.Errors():
		err = e.Err
		cb(&ProducerMessage{e.Msg}, e.Err)
	}

	return
}
