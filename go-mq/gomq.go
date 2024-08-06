package gomq

import "sync"

type BlockingQueue struct {
	queue chan interface{}
	wg    sync.WaitGroup
}

// NewBlockingQueue 创建一个新的阻塞队列
func NewBlockingQueue(size int) *BlockingQueue {
	return &BlockingQueue{
		queue: make(chan interface{}, size),
	}
}

// Enqueue 将元素添加到队列
func (bq *BlockingQueue) Enqueue(item interface{}) {
	bq.wg.Add(1)
	bq.queue <- item
}

// Dequeue 从队列中移除并返回一个元素
func (bq *BlockingQueue) Dequeue() interface{} {
	item := <-bq.queue
	bq.wg.Done()
	return item
}

// Size 返回队列中的元素数量
func (bq *BlockingQueue) Size() int {
	return len(bq.queue)
}

// Wait 等待所有元素被处理
func (bq *BlockingQueue) Wait() {
	bq.wg.Wait()
}
