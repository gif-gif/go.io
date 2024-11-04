package main

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
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
