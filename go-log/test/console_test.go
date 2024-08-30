package test

import (
	"github.com/gif-gif/go.io/go-log"
	"log"
	"testing"
)

func TestNewConsoleLog(t *testing.T) {
	l := golog.NewConsoleLog()

	l.WithHook(func(msg *golog.Message) {
		log.Println("this is hook", msg)
	})

	l.Debug("hi goio", "aa", 100)
	l.DebugF("hi %s", "goio")
	l.WithTag("u1").Debug("hi goio")
	l.WithTag("u1", "u-1").Warn("hi goio")
	l.WithTag("u1").WithField("name", "goio").Error("hi goio")
	l.WithTag("u1").WithField("id", 101).Panic("hi goio")
	l.WithTag("u1").WithField("id", 101).Fatal("hi goio")
}
