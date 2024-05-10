package golog

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewFileLog(t *testing.T) {
	l := NewFileLog(
		FilePathOption("logs/"),
		FileMaxSizeOption(1<<20),
	)

	l.WithHook(func(msg *Message) {
		log.Println("this is hook", msg)
	})

	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			l.Debug("hi hnatao", "aa", 100)
			l.DebugF("hi %s", "hnatao")
			l.WithTag("u1").Debug("hi hnatao")
			l.WithTag("u1", "u-1").Warn("hi hnatao")
			l.WithTag("u1").WithField("name", "hnatao").Error("hi hnatao")
			l.WithTag("u1").WithField("id", 101).Panic("hi hnatao")
			//l.WithTag("u1").WithField("id", 101).Fatal("hi hnatao")
		}()
	}

	wg.Wait()

	time.Sleep(3 * time.Second)
}
