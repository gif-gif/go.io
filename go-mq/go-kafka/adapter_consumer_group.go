package gokafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	gocontext "github.com/gif-gif/go.io/go-context"
	goutils "github.com/gif-gif/go.io/go-utils"
)

// 分组
type group struct {
	*GoKafka
	id      string
	handler ConsumerHandler
}

func (g group) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (group) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g group) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return fmt.Errorf("关闭会话上下文: %s", session.Context().Err())

		case msg, ok := <-claim.Messages():
			if !ok {
				return fmt.Errorf("消费通道关闭: groupId=%s topic=%s partition=%d", g.id, claim.Topic(), claim.Partition())
			}
			g.doHandler(msg, session)
		}
	}
}

func (g group) doHandler(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) (err error) {
	// 消息key
	key := string(msg.Key)

	// 消息试题
	m := goutils.M{
		"topic":     msg.Topic,
		"key":       key,
		"partition": msg.Partition,
		"offset":    msg.Offset,
		"timestamp": msg.Timestamp.Format("2006-01-02 15:04:05"),
	}

	// 填充数据
	{
		// body
		if len(msg.Value) > 0 {
			var body interface{}
			if err = json.Unmarshal(msg.Value, &body); err == nil {
				m["body"] = body
			} else {
				m["body"] = string(msg.Value)
			}
		}

		// headers
		for _, i := range msg.Headers {
			var headers = map[string]string{}
			headers[string(i.Key)] = string(i.Value)
			m["headers"] = headers
		}
	}

	// 定义上下文，关联到 session context 以便能够取消
	ctx := gocontext.WithParent(session.Context()).WithLog()
	ctx.Log.WithTag("gokafka-consumer-group", g.id).WithField("msg", m)

	if g.redis == nil {
		ctx.Log.WithField("msg", m)
		g.handleMsg(ctx, msg, session)
		return nil
	}
	// uniq key
	{
		var uniqKey string
		if key != "" {
			uniqKey = fmt.Sprintf("%s:%s", g.id, key)
		} else {
			uniqKey = fmt.Sprintf("%s:%s:%s", g.id, msg.Topic, goutils.MD5([]byte(g.id+msg.Topic+string(msg.Value))))
		}
		if g.redis != nil {
			ok := g.redis.SetNX(uniqKey, goutils.M{
				"topic":     msg.Topic,
				"body":      m["body"],
				"headers":   m["headers"],
				"timestamp": m["timestamp"],
			}.String(), 300*time.Second).Val()
			if !ok {
				ctx.Log.Warn("消息消费失败，并发消费")
				return
			}
			defer func() {
				g.redis.Del(uniqKey)
			}()
		}
	}

	// 建立缓存
	if g.redis != nil && key != "" {
		g.redis.Set1(key, goutils.M{
			"topic":     msg.Topic,
			"body":      m["body"],
			"headers":   m["headers"],
			"timestamp": m["timestamp"],
		}.String(), time.Hour)
	}

	// 打印日志
	t1 := time.Now()
	defer func() {
		ctx.Log.WithField("执行时间", fmt.Sprintf("%f", float64(time.Now().Sub(t1).Milliseconds())/1e3))
		if err != nil {
			ctx.Log.Error("消息消费失败", err)
			return
		}
		//ctx.Log.Debug("消息消费成功")
	}()

	// 执行业务方法
	g.handleMsg(ctx, msg, session)

	// 删除缓存
	if g.redis != nil && key != "" {
		g.redis.Del(key)
	}

	return
}

func (g group) handleMsg(ctx *gocontext.Context, msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
	// 执行业务方法
	if err := g.handler(ctx, &ConsumerMessage{ConsumerMessage: msg, GroupSession: session}, nil); err != nil {
		return
	}

	// 提交
	session.MarkMessage(msg, "")
}
