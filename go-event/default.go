package goevent

import "sync"

var (
	__event *Event
	__once  sync.Once
)

func Default() *Event {
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
