package goevent

import "sync"

var (
	__event *GoEvent
	__once  sync.Once
)

func Default() *GoEvent {
	return __event
}

func Publish(topic string, data interface{}) {
	__once.Do(func() {
		__event = New()
	})

	__event.Publish(topic, data)
}

func Subscribe(topic string, fn SubscribeFunc) {
	__once.Do(func() {
		__event = New()
	})

	__event.Subscribe(topic, fn)
}
