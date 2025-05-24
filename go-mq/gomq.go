package gomq

import (
	"errors"
	"sync"
	"time"
)

type BlockingQueue struct {
	queue    chan interface{}
	wg       sync.WaitGroup
	closed   bool
	closeMux sync.RWMutex
}

// NewBlockingQueue 创建一个新的阻塞队列
func NewBlockingQueue(size int) *BlockingQueue {
	return &BlockingQueue{
		queue:  make(chan interface{}, size),
		closed: false,
	}
}

// Enqueue 将元素添加到队列
func (bq *BlockingQueue) Enqueue(item interface{}) error {
	bq.closeMux.RLock()
	if bq.closed {
		bq.closeMux.RUnlock()
		return errors.New("queue is closed")
	}

	bq.closeMux.RUnlock()
	// 使用 select 防止死锁
	select {
	case bq.queue <- item:
		// 如果成功添加元素，增加等待组计数
		bq.wg.Add(1)
		return nil
	default:
		return errors.New("queue is full")
	}
}

// EnqueueBlocking 阻塞式将元素添加到队列
func (bq *BlockingQueue) EnqueueBlocking(item interface{}) error {
	bq.closeMux.RLock()
	if bq.closed {
		bq.closeMux.RUnlock()
		return errors.New("queue is closed")
	}

	bq.wg.Add(1)
	bq.closeMux.RUnlock()

	// 此处会阻塞直到队列有空间
	select {
	case bq.queue <- item:
		return nil
	}
}

// Dequeue 从队列中移除并返回一个元素，如果队列为空则阻塞
func (bq *BlockingQueue) Dequeue() (interface{}, error) {
	item, ok := <-bq.queue
	if !ok {
		return nil, errors.New("queue is closed")
	}
	bq.wg.Done()
	return item, nil
}

// DequeueWithTimeout 带超时的出队操作
func (bq *BlockingQueue) DequeueWithTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case item, ok := <-bq.queue:
		if !ok {
			return nil, errors.New("queue is closed")
		}
		bq.wg.Done()
		return item, nil
	case <-time.After(timeout):
		return nil, errors.New("dequeue timeout")
	}
}

// Size 返回队列中的元素数量（注意：这只是一个瞬时值）
func (bq *BlockingQueue) Size() int {
	return len(bq.queue)
}

// Wait 等待所有元素被处理
func (bq *BlockingQueue) Wait() {
	bq.wg.Wait()
}

// Close 关闭队列，不再接受新元素
func (bq *BlockingQueue) Close() {
	bq.closeMux.Lock()
	defer bq.closeMux.Unlock()

	if !bq.closed {
		bq.closed = true
		close(bq.queue)
	}
}

// IsClosed 检查队列是否已关闭
func (bq *BlockingQueue) IsClosed() bool {
	bq.closeMux.RLock()
	defer bq.closeMux.RUnlock()
	return bq.closed
}
