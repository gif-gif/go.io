package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	goevent "github.com/gif-gif/go.io/go-event"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"github.com/gogf/gf/util/gconv"
	"time"
)

type EventField1 struct {
	A string
}
type EventField2 struct {
	B string
}

type TestEvent struct {
	A *EventField1
	B *EventField2
}

func main() {
	goio.Init(goio.DEVELOPMENT)
	simpleTest()
	<-gocontext.Cancel().Done()
}

func simpleTest() {
	event := goevent.New()
	event.Subscribe("test", func(msg goevent.Message) {
		golog.WithTag("goevent").Info(msg)
	})
	event.Publish("test", "test")
	time.Sleep(1 * time.Second)
	event.UnSubscribe("test")
}

func channelSizeTest() {
	event := goevent.New(100)
	event.Subscribe("test", func(msg goevent.Message) {
		golog.WithTag("goevent").Info(msg)
	})
	for i := 0; i < 100; i++ {
		event.Publish("test", "test-"+gconv.String(i))
	}
}
