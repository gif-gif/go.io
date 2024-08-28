package golog

import (
	"log"
	"testing"
)

func TestDefault(t *testing.T) {
	//Default().SetAdapter(&FileAdapter{})

	WithHook(func(msg *Message) {
		log.Println("this is hook", msg)
	})

	Debug("hi goio", "aa", 100)
	DebugF("hi %s", "goio")
	WithTag("u1").Debug("hi goio")
	WithTag("u1", "u-1").Debug("hi goio")
	WithTag("u1").WithField("name", "goio").Debug("hi goio")
	WithTag("u1").WithField("id", 101).Debug("hi goio")
	WithTag("u1").WithField("id", 101).Debug("hi goio")
}
