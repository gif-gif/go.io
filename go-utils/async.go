package goutils

import (
	"context"
	golog "github.com/gif-gif/go.io/go-log"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

type ErrorGroup struct {
	maxWorkers int
	errGroup   *errgroup.Group
	ctx        context.Context
}

// 当并发执行过程中有错误时，会自动取消其他所有任务,通过CancelContext 取消来实现
// （最大并发数为 maxWorkers,超过阻塞等待 ）
func NewErrorGroup(context context.Context, maxWorkers int) ErrorGroup {
	e := ErrorGroup{}
	e.maxWorkers = maxWorkers
	e.errGroup, e.ctx = errgroup.WithContext(context)
	e.errGroup.SetLimit(e.maxWorkers)
	return e
}

func (e *ErrorGroup) Submit(fn ...func() error) {
	for _, f := range fn {
		e.errGroup.Go(f)
	}
}

func (e *ErrorGroup) TryGo(fn func() error) bool {
	return e.errGroup.TryGo(fn)
}

func (e *ErrorGroup) Wait() error {
	return e.errGroup.Wait()
}

func (e *ErrorGroup) IsContextDone() bool {
	return IsContextDone(e.ctx)
}

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

// 异步执行（安全） RunSafe
func AsyncFunc(fn func()) {
	go func() {
		defer Recovery(nil)
		fn()
	}()
}

// 异步并发执行（安全), 建议使用NewErrorGroup 替代
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
//
// 建议使用NewErrorGroup 替代
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

// 返回函数执行时间
func MeasureExecutionTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
	}
	return false
}
