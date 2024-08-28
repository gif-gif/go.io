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

			l.Debug("hi goio", "aa", 100)
			l.DebugF("hi %s", "goio")
			l.WithTag("u1").Debug("hi goio")
			l.WithTag("u1", "u-1").Warn("hi goio")
			l.WithTag("u1").WithField("name", "goio").Error("hi goio")
			l.WithTag("u1").WithField("id", 101).Panic("hi goio")
			//l.WithTag("u1").WithField("id", 101).Fatal("hi goio")
		}()
	}

	wg.Wait()

	time.Sleep(3 * time.Second)
}
