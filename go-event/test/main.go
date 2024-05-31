package main

import (
	goevent "github.com/gif-gif/go.io/go-event"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init()
	event := goevent.New()
	event.Subscribe("test", func(msg goevent.Message) {
		golog.WithTag("goevent").Info(msg)
	})
	event.Publish("test", "test")
	time.Sleep(time.Duration(1) * time.Second)
}
