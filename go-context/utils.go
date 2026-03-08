package gocontext

import (
	"context"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"
)

func CatchErrorAndThrow() {
	if err := recover(); err != nil {
		log.Printf("error:%v stack:%v", err, string(debug.Stack()))
		log.Panic(err)
	}
}

func RoutineWrapper(fn func()) {
	defer CatchErrorAndThrow()
	fn()
}

func RoutineWrapper1[T1 any](fn func(T1), t1 T1) {
	defer CatchErrorAndThrow()
	fn(t1)
}

func CreateContextWithCancel(parentCtx context.Context, timeout time.Duration, stopChans ...chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parentCtx)

	if timeout > 0 {
		go RoutineWrapper(func() {
			timer := time.NewTimer(timeout)
			defer timer.Stop()
			select {
			case <-timer.C:
				cancel()
				break
			case <-ctx.Done():
				break
			}
		})
	}

	for _, stopChan := range stopChans {
		if stopChan == nil {
			continue
		}
		go RoutineWrapper(func() {
			select {
			case <-stopChan:
				cancel()
				break
			case <-ctx.Done():
				break
			}
		})
	}

	return ctx, cancel
}

// 多渠道合并（Channel Merging） 模式，高并发微服务中。
// 它的核心功能是将多个 stopChan 合并为一个，只要其中任何一个关闭，或者手动调用返回的 cancelFn，最终的 outChan 就会关闭。
func MergeStopChans(stopChans ...chan struct{}) (chan struct{}, context.CancelFunc) {
	outChan := make(chan struct{})
	outChanClosed := int32(0)

	cancelFn := func() {
		if atomic.SwapInt32(&outChanClosed, 1) == 0 {
			close(outChan)
		}
	}

	for _, stopChan := range stopChans {
		if stopChan == nil {
			continue
		}
		go RoutineWrapper1(func(c chan struct{}) {
			select {
			case <-stopChan:
				cancelFn()
			case <-outChan:
			}
		}, stopChan)
	}

	return outChan, cancelFn
}
