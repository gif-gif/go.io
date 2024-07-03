# 定时任务
```go
 需要支持荥
 带参数的任务函数 定时执行还没支持
 每个月某一天几点执行

```
```
package main

import (
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"time"
)

func main() {
	DataChan := make(chan []byte, 20)

	n := 1
	cron, err := gojob.New()
	if err != nil {
		golog.WithTag("gojob").Error(err)
	}
	defer cron.Stop()
	defer close(DataChan)

	cron.Start()
	cron.Second(func() {
		if r := recover(); r != nil {
			golog.Error(r)
		}

		golog.WithTag("gojob").Info("testing")
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
				golog.WithTag("gojob").Info(string(data))
			}
		}
	}()

	time.Sleep(time.Second * 5)

}

```