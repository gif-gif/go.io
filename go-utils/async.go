package goutils

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"sync"
	"time"
)

// 捕获panic
func Recovery(errFn func(err any)) {
	if r := recover(); r != nil {
		if errFn != nil {
			errFn(r)
		} else {
			golog.Error(r)
		}
	}
}

// 异步执行（安全）errFn = nil 时自动Recovery 不会Panic
func AsyncFuncPanic(fn func(), errFn func(err any)) {
	go func() {
		defer Recovery(errFn)
		fn()
	}()
}

// 异步执行（安全）
func AsyncFunc(fn func()) {
	go func() {
		defer Recovery(nil)
		fn()
	}()
}

// 异步并发执行（安全
func AsyncFuncGroup(fns ...func()) {
	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)
		func(fn func()) {
			AsyncFunc(func() {
				defer wg.Done()
				fn()
			})
		}(fn)
	}

	wg.Wait()
}

// 异步并发执行（安全）errFn = nil 时自动Recovery 不会Panic
func AsyncFuncGroupPanic(errFn func(err any), fns ...func()) {
	var wg sync.WaitGroup
	for _, fn := range fns {
		wg.Add(1)
		func(fn func()) {
			AsyncFuncPanic(func() {
				defer wg.Done()
				fn()
			}, errFn)
		}(fn)
	}

	wg.Wait()
}

// 异步并发执行（安全) 只要其中一个执行完成（不能嵌套使用），释放所有等待锁。目的是比执行速度
func AsyncFuncGroupOneSuccess(fns ...func()) {
	var wg sync.WaitGroup
	lockCount := len(fns)
	for i := 0; i < lockCount; i++ {
		wg.Add(1)
	}

	for _, fn := range fns {
		func(fn func()) {
			AsyncFunc(func() {
				defer wg.Done()
				fn()
				//有一个执行完成后其他的都退出阻塞
				for i := 0; i < lockCount-1; i++ {
					wg.Done()
				}
			})
		}(fn)
	}

	wg.Wait()
}

func testRaceSpeed() {
	var fns []func()
	fns = append(fns, func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Hello 5")
	})

	fns = append(fns, func() {
		time.Sleep(10 * time.Second)
		fmt.Println("Hello 1")
	})

	fns = append(fns, func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Hello 3")
	})

	AsyncFuncGroupOneSuccess(fns...)
}
