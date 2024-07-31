package main

import (
	"fmt"
	golock "github.com/gif-gif/go.io/go-lock"
	gopool "github.com/gif-gif/go.io/go-pool"
	"github.com/gif-gif/go.io/goio"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	testSyncLock()
}

func testSyncLock() {
	var (
		count int
	)
	lock := golock.New()
	//并发池
	gp := gopool.NewFixedSizePool(10, 10) // 可以并发10个
	defer gp.StopAndWait()
	group := gp.NewTaskGroup()
	for i := 0; i < 2; i++ {
		group.Submit(func() {
			for i := 100000; i > 0; i-- {
				lock.WLockFunc(func(parameters ...any) {
					count++
				})
			}
		})
	}

	group.Wait()
	fmt.Println(count) //输出 200000
}
