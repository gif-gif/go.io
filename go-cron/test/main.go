package main

import (
	gocrons "github.com/gif-gif/go.io/go-cron"
	golog "github.com/gif-gif/go.io/go-log"
	"time"
)

func main() {
	DataChan := make(chan []byte, 20)
	n := 1
	cron := gocrons.New()
	defer cron.Stop()
	defer close(DataChan)

	cron.Start()
	cron.Second(func() {
		if r := recover(); r != nil {
			golog.Error(r)
		}

		golog.WithTag("gocrons").Info("testing")
		n++
		if n > 5 {
			n = 0
			cron.Stop()
		}
		DataChan <- []byte("json")
	})

	go func() {
		for {
			select {
			case data := <-DataChan:
				golog.WithTag("gocrons").Info(string(data))
			}
		}
	}()

	time.Sleep(time.Second * 5)

}
