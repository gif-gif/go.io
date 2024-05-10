package golog

import (
	"log"
	"testing"
)

func TestNewConsoleLog(t *testing.T) {
	l := NewConsoleLog()

	l.WithHook(func(msg *Message) {
		log.Println("this is hook", msg)
	})

	l.Debug("hi hnatao", "aa", 100)
	l.DebugF("hi %s", "hnatao")
	l.WithTag("u1").Debug("hi hnatao")
	l.WithTag("u1", "u-1").Warn("hi hnatao")
	l.WithTag("u1").WithField("name", "hnatao").Error("hi hnatao")
	l.WithTag("u1").WithField("id", 101).Panic("hi hnatao")
	l.WithTag("u1").WithField("id", 101).Fatal("hi hnatao")
}
