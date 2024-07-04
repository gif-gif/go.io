package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	goevent "github.com/gif-gif/go.io/go-event"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
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
	event := goevent.New()
	event.Subscribe("test", func(msg goevent.Message) {
		golog.WithTag("goevent").Info(msg)
	})
	event.Publish("test", "test")

	event.Subscribe("testFields", func(msg goevent.Message) {
		t := msg.Data.(*TestEvent)
		golog.WithTag("goevent").InfoF("eventData:%+v", t)
	})

	event.Publish("testFields", &TestEvent{
		A: &EventField1{
			A: "AAAAAA",
		},
		B: &EventField2{B: "BBBBB"},
	})

	time.Sleep(10 * time.Second)

	<-gocontext.Cancel().Done()
}
