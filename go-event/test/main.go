package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	goevent "github.com/gif-gif/go.io/go-event"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	event := goevent.New()
	event.Subscribe("test", func(msg goevent.Message) {
		golog.WithTag("goevent").Info(msg)
	})
	event.Publish("test", "test")
	<-gocontext.Cancel().Done()
}
