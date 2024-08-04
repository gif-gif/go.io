package goevent

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"sync"
)

type GoEvent struct {
	subscribes map[string][]MessageChan
	mu         sync.RWMutex
}

func New() *GoEvent {
	return &GoEvent{subscribes: map[string][]MessageChan{}}
}

// 发布 执行当前topic 对应的所有订阅者
func (ev *GoEvent) Publish(topic string, data interface{}) {
	ev.mu.RLock()
	defer ev.mu.RUnlock()

	if chs, ok := ev.subscribes[topic]; ok {
		channels := append([]MessageChan{}, chs...)
		goutils.AsyncFunc(func() {
			for _, ch := range channels {
				ch <- Message{Topic: topic, Data: data}
			}
		})
	}
}

// 订阅：一个topic可以对应多个处理器，（topic->handler 的关系是1:n）,一次添加一个订阅者
func (ev *GoEvent) Subscribe(topic string, fn SubscribeFunc) {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	if _, ok := ev.subscribes[topic]; !ok {
		ev.subscribes[topic] = []MessageChan{}
	}

	ch := make(chan Message)
	ev.subscribes[topic] = append(ev.subscribes[topic], ch)

	goutils.AsyncFunc(func() {
		for {
			select {
			case msg := <-ch:
				goutils.AsyncFunc(func() {
					fn(msg)
				})
			}
		}
	})
}

func (ev *GoEvent) UnSubscribe(topic string) {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	if chs, ok := ev.subscribes[topic]; ok {
		channels := append([]MessageChan{}, chs...)
		for _, ch := range channels {
			close(ch)
		}

		delete(ev.subscribes, topic)
	}
}
