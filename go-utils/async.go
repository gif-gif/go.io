package goutils

import (
	"context"
	"fmt"
	"sync"
	"time"

	golog "github.com/gif-gif/go.io/go-log"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
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

func (e *ErrorGroup) GetContext() context.Context {
	return e.ctx
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
			logx.Errorf("AsyncFunc recover: %v", r)
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

// --- Safe Goroutine ---

// CatchErrorAndThrow recovers from panics (use with defer).
func CatchErrorAndThrow() {
	if r := recover(); r != nil {
		panic(r)
	}
}

// Go runs fn in a goroutine with panic recovery.
func Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Safe Goroutine.Go() panic: %v\n", r)
			}
		}()
		fn()
	}()
}

// Go1 runs fn(t1) in a goroutine with panic recovery.
func Go1[T1 any](fn func(T1), t1 T1) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Safe Goroutine.Go1() panic: %v\n", r)
			}
		}()
		fn(t1)
	}()
}

// Go2 runs fn(t1, t2) in a goroutine with panic recovery.
func Go2[T1, T2 any](fn func(T1, T2), t1 T1, t2 T2) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Safe Goroutine.Go2() panic: %v\n", r)
			}
		}()
		fn(t1, t2)
	}()
}

// Go3 runs fn(t1, t2, t3) in a goroutine with panic recovery.
func Go3[T1, T2, T3 any](fn func(T1, T2, T3), t1 T1, t2 T2, t3 T3) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Safe Goroutine.Go3() panic: %v\n", r)
			}
		}()
		fn(t1, t2, t3)
	}()
}
