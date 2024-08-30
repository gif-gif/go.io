package test

import (
	"github.com/gif-gif/go.io/go-log"
	"log"
	"testing"
)

func TestDefault(t *testing.T) {
	//Default().SetAdapter(&FileAdapter{})

	golog.WithHook(func(msg *golog.Message) {
		log.Println("this is hook", msg)
	})

	golog.Debug("hi goio", "aa", 100)
	golog.DebugF("hi %s", "goio")
	golog.WithTag("u1").Debug("hi goio")
	golog.WithTag("u1", "u-1").Debug("hi goio")
	golog.WithTag("u1").WithField("name", "goio").Debug("hi goio")
	golog.WithTag("u1").WithField("id", 101).Debug("hi goio")
	golog.WithTag("u1").WithField("id", 101).Debug("hi goio")
}
