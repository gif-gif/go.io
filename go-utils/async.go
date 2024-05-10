package goutils

import (
	golog "github.com/jiriyao/go.io/go-log"
	"sync"
)

// 捕获panic
func Recovery() {
	if r := recover(); r != nil {
		golog.Error(r)
	}
}

// 异步执行（安全）
func AsyncFunc(fn func()) {
	go func() {
		defer Recovery()
		fn()
	}()
}

// 异步并发执行（安全）
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
