package main

import (
	"context"
	"fmt"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	gopool "github.com/gif-gif/go.io/go-pool"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func timeCost(start time.Time) {
	tc := time.Since(start)
	fmt.Printf("time cost = %v\n", tc.Seconds())
}

// TODO: 并发怎么迅速取消或者停止，http 比如请求最快的直接返回
func main() {
	goio.Init(goio.DEVELOPMENT)
	defer timeCost(time.Now())
	testStopPool()
}

func testDynamicSize() {
	gp := gopool.NewDynamicSizePool(100, 10)
	defer gp.StopAndWait()

	cron, _ := gojob.New()
	defer cron.Stop()
	cron.Start()
	cron.SecondX(nil, 1, func() {
		gp.PrintPoolStats()
	})

	for i := 0; i < 1000; i++ {
		n := i
		gp.Submit(func() {
			fmt.Printf("Running task #%d\n", n)
			time.Sleep(1 * time.Second)
		})
	}

	golog.InfoF("end of Submit")
}

func testFixedSize() {
	gp := gopool.NewFixedSizePool(100, 10)
	defer gp.StopAndWait()

	cron, _ := gojob.New()
	defer cron.Stop()
	cron.Start()
	cron.SecondX(nil, 1, func() {
		gp.PrintPoolStats()
	})

	for i := 0; i < 1000; i++ {
		n := i
		gp.Submit(func() {
			fmt.Printf("Running task #%d\n", n)
			time.Sleep(1 * time.Second)
		})
	}

	golog.InfoF("end of Submit")
}

func testContext() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10) //超时或者取消ctx 时池子会被关闭，未开始执行的任务会被取消执行

	gp := gopool.NewContextPool(10, 10, ctx)
	defer gp.StopAndWait()

	cron, _ := gojob.New()
	defer cron.Stop()
	cron.Start()
	cron.SecondX(nil, 1, func() {
		gp.PrintPoolStats()
	})

	for i := 0; i < 1000; i++ {
		n := i
		gp.Submit(func() {
			fmt.Printf("Task #%d started\n", n)
			time.Sleep(1 * time.Second)
			fmt.Printf("Task #%d finished\n", n)
		})
	}

	golog.InfoF("end of Submit")
}

func testTaskGroup() {
	gp := gopool.NewDynamicSizePool(100, 1000)
	defer gp.StopAndWait()

	cron, _ := gojob.New()
	defer cron.Stop()
	cron.Start()
	cron.SecondX(nil, 1, func() {
		gp.PrintPoolStats()
	})

	group := gp.NewTaskGroup()

	for i := 0; i < 1000; i++ {
		n := i
		group.Submit(func() {
			fmt.Printf("Task #%d started\n", n)
			time.Sleep(1 * time.Second)
			fmt.Printf("Task #%d finished\n", n)
		})
	}

	group.Wait() // wait for tasks to finish
	golog.InfoF("end of TaskGroup")
}

func testGroupContext() {
	gp := gopool.NewDynamicSizePool(10, 1000)
	defer gp.StopAndWait()

	cron, _ := gojob.New()
	defer cron.Stop()
	cron.Start()
	cron.SecondX(nil, 1, func() {
		gp.PrintPoolStats()
	})

	group, _ := gp.NewGroupContext() //可以用 ctx 和 group 组合使用

	for i := 0; i < 1000; i++ {
		n := i
		group.Submit(func() error {
			fmt.Printf("Task #%d started\n", n)
			time.Sleep(1 * time.Second)
			fmt.Printf("Task #%d finished\n", n)

			if n > 1 { //出错后其他未开始任务执行会被取消
				return fmt.Errorf("test group error")
			}

			return nil
		})
	}

	err := group.Wait() // wait for tasks to finish
	if err != nil {
		golog.InfoF("end of GroupContext", err)
	} else {
		golog.InfoF("end of GroupContext")
	}
}

func testStopPool() {
	gp := gopool.NewDynamicSizePool(1, 100)
	defer gp.StopAndWait()
	group := gp.NewTaskGroup() //可以用 ctx 和 group 组合使用
	var ctx context.Context
	for i := 0; i < 50; i++ {
		n := i
		group.Submit(func() {
			if gp.Stopped() {
				return
			}
			fmt.Printf("Task #%d started\n", n)
			time.Sleep(1 * time.Second)
			fmt.Printf("Task #%d finished\n", n)
			if n == 2 {
				ctx = gp.Stop()
				goutils.AsyncFunc(func() {
					if ctx != nil {
						select {
						case <-ctx.Done():
							golog.Error("Worker pool did not stop gracefully", gp.Stopped())
						}
					}
				})
			}
		})
	}

	group.Wait() // wait for tasks to finish but Worker pool did not stop gracefully
	golog.InfoF("Worker pool finished")
}
