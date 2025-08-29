package gopool

import (
	"github.com/panjf2000/ants/v2"
	"time"
)

type GoPool struct {
	pool *ants.Pool
}

type PoolStat struct {
	RunningWorkers  int64
	IdleWorkers     int
	SubmittedTasks  uint64
	WaitingTasks    uint64
	SuccessfulTasks uint64
	FailedTasks     uint64
	CompletedTasks  uint64
}

//	ants will pre-malloc the whole capacity of pool when calling ants.NewPool.
//
// p, _ := ants.NewPool(100000, ants.WithPreAlloc(true))
func New(maxWorkers int, options ...ants.Option) (*GoPool, error) {
	p, err := ants.NewPool(maxWorkers, options...)
	return &GoPool{pool: p}, err
}

func (g *GoPool) GetPool() *ants.Pool {
	return g.pool
}

func (g *GoPool) Submit(fn func()) error {
	return g.pool.Submit(fn)
}

// 释放关闭此池并释放工作队列。
func (g *GoPool) Release() {
	g.pool.Release()
}

// 只要调用 Reboot() 方法，就可以重新激活一个之前已经被销毁掉的池，并且投入使用。
func (g *GoPool) Reboot() {
	g.pool.Reboot()
}

// ReleaseTimeout就像Release，但带有超时，等待所有工作者退出后再超时。
func (g *GoPool) ReleaseTimeout(timeout time.Duration) error {
	return g.pool.ReleaseTimeout(timeout)
}

// 可用Worker数量，-1 表示无限制
func (g *GoPool) Free() {
	g.pool.Free()
}

// 调整更改了此池的容量，请注意，这对无限池或预分配池是无效的。线程安全
func (g *GoPool) Tune(size int) {
	g.pool.Tune(size)
}

func (g *GoPool) IsClosed() bool {
	return g.pool.IsClosed()
}

// 运行中的Worker数量
func (g *GoPool) Running() int {
	return g.pool.Running()
}

// 等待返回等待执行的任务数量。
func (g *GoPool) Waiting() int {
	return g.pool.Waiting()
}

// Cap返回该池的容量
func (g *GoPool) Cap() int {
	return g.pool.Cap()
}

//
//func (g *GoPool) PoolStats() PoolStat {
//	return PoolStat{
//		RunningWorkers:  g.pool.RunningWorkers(),
//		SubmittedTasks:  g.pool.SubmittedTasks(),
//		WaitingTasks:    g.pool.WaitingTasks(),
//		SuccessfulTasks: g.pool.SuccessfulTasks(),
//		FailedTasks:     g.pool.FailedTasks(),
//		CompletedTasks:  g.pool.CompletedTasks(),
//	}
//}
//
//func (g *GoPool) PrintPoolStats() {
//	ps := g.PoolStats()
//	golog.WithTag("GoPool").InfoF("RunningWorkers: %d", ps.RunningWorkers)
//	golog.WithTag("GoPool").InfoF("IdleWorkers: %d", ps.IdleWorkers)
//	golog.WithTag("GoPool").InfoF("SubmittedTasks: %d", ps.SubmittedTasks)
//	golog.WithTag("GoPool").InfoF("WaitingTasks: %d", ps.WaitingTasks)
//	golog.WithTag("GoPool").InfoF("SuccessfulTasks: %d", ps.SuccessfulTasks)
//	golog.WithTag("GoPool").InfoF("FailedTasks: %d", ps.FailedTasks)
//	golog.WithTag("GoPool").InfoF("CompletedTasks: %d", ps.CompletedTasks)
//	golog.WithTag("GoPool").InfoF("----------------------------------------------------------------")
//}
