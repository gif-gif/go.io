package main

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	gopool "github.com/gif-gif/go.io/go-pool"
	"github.com/gif-gif/go.io/goio"
	"github.com/panjf2000/ants/v2"
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
	pool, _ := gopool.New(100, ants.WithPreAlloc(true))
	err := pool.Submit(func() {

	})
	if err != nil {
		golog.ErrorF("Submit failed: %v", err)
	}
}

func testDynamicSize() {
	golog.InfoF("end of Submit")
}

func testFixedSize() {
	golog.InfoF("end of Submit")
}

func testContext() {
	golog.InfoF("end of Submit")
}

func testTaskGroup() {
	golog.InfoF("end of TaskGroup")
}

func testGroupContext() {

}

func testStopPool() {
	golog.InfoF("Worker pool finished")
}
