package goevent

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"sync"
)

// 支持同步和异步执行(发送和接收)
//
// 默认是同步执行
type GoEvent struct {
	subscribes   map[string][]MessageChan
	channelSize  int // 消息通道同时处理大小，默认没有限制
	mu           sync.RWMutex
	DefaultTopic string
}

func New(channelSize ...int) *GoEvent {
	if len(channelSize) == 0 {
		return &GoEvent{subscribes: map[string][]MessageChan{}, channelSize: 0}
	} else {
		return &GoEvent{subscribes: map[string][]MessageChan{}, channelSize: channelSize[0]}
	}
}

// 发布 执行当前topic 对应的所有订阅者, async=true 则异步执行(并发执行无序)，否则同步执行保证channel发送顺序
func (ev *GoEvent) Publish(topic string, data interface{}, async ...bool) {
	ev.mu.RLock()
	defer ev.mu.RUnlock()

	if chs, ok := ev.subscribes[topic]; ok {
		channels := append([]MessageChan{}, chs...)
		if len(async) > 0 && async[0] { //并发执行
			goutils.AsyncFunc(func() {
				for _, ch := range channels {
					ch <- Message{Topic: topic, Data: data}
				}
			})
		} else {
			for _, ch := range channels {
				ch <- Message{Topic: topic, Data: data}
			}
		}
	}
}

// 订阅：一个topic可以对应多个处理器，（topic->handler 的关系是1:n）,一次添加一个订阅者
func (ev *GoEvent) Subscribe(topic string, fn SubscribeFunc, async ...bool) {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	if _, ok := ev.subscribes[topic]; !ok {
		ev.subscribes[topic] = []MessageChan{}
	}
	var ch chan Message

	if ev.channelSize == 0 {
		ch = make(chan Message)
	} else {
		ch = make(chan Message, ev.channelSize)
	}
	if ev.DefaultTopic == "" {
		ev.DefaultTopic = topic
	}
	ev.subscribes[topic] = append(ev.subscribes[topic], ch)

	goutils.AsyncFunc(func() {
		if len(async) > 0 && async[0] { //并发执行
			for msg := range ch {
				goutils.AsyncFunc(func() {
					fn(msg)
				})
			}
		} else {
			for msg := range ch {
				fn(msg)
			}
		}
	})
}

// 取消订阅(topic 对应的所有订阅者)
func (ev *GoEvent) UnSubscribe(topic string) {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	if chs, ok := ev.subscribes[topic]; ok {
		channels := append([]MessageChan{}, chs...)
		for _, ch := range channels {
			goutils.AsyncFunc(func() {
				close(ch)
			})
		}

		delete(ev.subscribes, topic)
	}
}

func (ev *GoEvent) UnSubscribeDefault() {
	ev.mu.Lock()
	defer ev.mu.Unlock()
	if chs, ok := ev.subscribes[ev.DefaultTopic]; ok {
		channels := append([]MessageChan{}, chs...)
		for _, ch := range channels {
			goutils.AsyncFunc(func() {
				close(ch)
			})
		}

		delete(ev.subscribes, ev.DefaultTopic)
	}
}
