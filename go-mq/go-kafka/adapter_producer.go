package gokafka

import (
	"errors"
	"github.com/IBM/sarama"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"time"
)

type producer struct {
	*GoKafka
	focus bool // 是否强制发送
}

func (p *producer) Client() sarama.Client {
	return p.GoKafka.Client
}

// 发送消息 - 同步- 批量
func (p *producer) SendMessages(msgs []IMessage) (err error) {
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
func (p *producer) SendMessage(msg IMessage) (partition int32, offset int64, err error) {
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

// 发送消息 - 异步 - 批量
func (p *producer) SendAsyncMessages(msgs []IMessage, cb MessageHandler) (err error) {
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

	producer, err := sarama.NewAsyncProducerFromClient(p.Client())
	if err != nil {
		return
	}
	defer producer.Close()
	for _, msg := range sendMsgs {
		producer.Input() <- msg
	}
	for i := 0; i < len(sendMsgs); i++ {
		select {
		case _msg := <-producer.Successes():
			cb(&ProducerMessage{_msg}, nil)
		case e := <-producer.Errors():
			err = e.Err
			cb(&ProducerMessage{e.Msg}, e.Err)
		}
	}

	return
}

// 发送消息 - 异步
func (p *producer) SendAsyncMessage(msg IMessage, cb MessageHandler) (err error) {
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
