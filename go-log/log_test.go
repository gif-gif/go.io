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

	Debug("hi hnatao", "aa", 100)
	DebugF("hi %s", "hnatao")
	WithTag("u1").Debug("hi hnatao")
	WithTag("u1", "u-1").Debug("hi hnatao")
	WithTag("u1").WithField("name", "hnatao").Debug("hi hnatao")
	WithTag("u1").WithField("id", 101).Debug("hi hnatao")
	WithTag("u1").WithField("id", 101).Debug("hi hnatao")
}
