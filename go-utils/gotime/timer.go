package gotime

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"sync"
	"time"
)

type Timer struct {
	duration  time.Duration
	stop      chan struct{}
	pause     chan struct{}
	resume    chan struct{}
	reset     chan time.Duration
	force     chan struct{}
	wg        sync.WaitGroup
	isPaused  bool
	isStopped bool
	mu        sync.Mutex
	ticker    *time.Ticker
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		duration: duration,
		// 使用buffered channels避免阻塞
		stop:      make(chan struct{}, 1),
		pause:     make(chan struct{}, 1),
		resume:    make(chan struct{}, 1),
		reset:     make(chan time.Duration, 1),
		force:     make(chan struct{}, 1),
		isStopped: true,
	}
}

// 不要重复调用Start方法
func (t *Timer) Start(task func()) {
	t.mu.Lock()
	// 重置停止状态
	t.isStopped = false
	t.mu.Unlock()

	t.wg.Add(1)
	if t.ticker == nil {
		t.ticker = time.NewTicker(t.duration)
	}

	goutils.AsyncFunc(func() {
		defer t.wg.Done()
		for {
			if t.ticker == nil {
				return
			}
			select {
			case <-t.ticker.C:
				t.mu.Lock()
				isPaused := t.isPaused
				t.mu.Unlock()
				if !isPaused {
					task()
				}
			case <-t.pause:
				t.mu.Lock()
				t.isPaused = true
				t.mu.Unlock()
			case <-t.resume:
				t.mu.Lock()
				t.isPaused = false
				t.mu.Unlock()
			case newDuration := <-t.reset:
				t.ticker.Reset(newDuration)
				t.duration = newDuration
			case <-t.force:
				task()
			case <-t.stop:
				return
			}
		}
	})
}

func (t *Timer) Stop() {
	t.mu.Lock()
	if !t.isStopped {
		t.isStopped = true
		if t.ticker != nil {
			t.ticker.Stop()
		}
		close(t.stop)
		// 重新创建stop channel为下次使用准备
		t.stop = make(chan struct{}, 1)
		t.ticker = nil
	}
	t.mu.Unlock()
	t.wg.Wait()
}

func (t *Timer) Pause() {
	t.mu.Lock()
	if !t.isStopped && !t.isPaused {
		select {
		case t.pause <- struct{}{}:
		default:
		}
	}
	t.mu.Unlock()
}

func (t *Timer) Resume() {
	t.mu.Lock()
	if !t.isStopped && t.isPaused {
		select {
		case t.resume <- struct{}{}:
		default:
		}
	}
	t.mu.Unlock()
}

func (t *Timer) ForceRun() {
	t.mu.Lock()
	if !t.isStopped {
		select {
		case t.force <- struct{}{}:
		default:
		}
	}
	t.mu.Unlock()
}

func (t *Timer) Reset(duration time.Duration) {
	t.mu.Lock()
	if !t.isStopped {
		select {
		case t.reset <- duration:
		default:
		}
	}
	t.mu.Unlock()
}
