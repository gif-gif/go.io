package golock

import "sync"

// RWMutex 也称为读写互斥锁，读写互斥锁就是读取/写入互相排斥的锁。它可以由任意数量的读取操作的 goroutine 或单个写入操作的 goroutine 持有。
// 读写互斥锁 RWMutex 类型有五个方法，Lock，Unlock，Rlock，RUnlock 和 RLocker。其中，RLocker 返回一个 Locker 接口，
// 该接口通过调用 rw.RLock 和 rw.RUnlock 来实现 Lock 和 Unlock 方法。
// 不能拷贝锁
type GoLock struct {
	Mutex  sync.Mutex   // 读锁时不能写，写锁时不能读取
	MuteRW sync.RWMutex //读写互斥锁，并发读取，单一写入。读多写少性能会好
}

func New() *GoLock {
	return &GoLock{}
}

// 共享内存加锁
func (g *GoLock) Lock(fn func(parameters ...any), parameters ...any) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	fn(parameters...)
}

// 加读锁
func (g *GoLock) GoRLock(fn func()) {
	g.MuteRW.RLock()
	defer g.MuteRW.RUnlock()
	fn()
}

// 加写锁
func (g *GoLock) GoWLock(fn func()) {
	g.MuteRW.Lock()
	defer g.MuteRW.Unlock()
	fn()
}
