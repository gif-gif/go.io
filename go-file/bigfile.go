package gofile

import (
	"context"
	"errors"
	"fmt"
	golock "github.com/gif-gif/go.io/go-lock"
	gomq "github.com/gif-gif/go.io/go-mq"
	gopool "github.com/gif-gif/go.io/go-pool"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"math"
	"os"
	"path/filepath"
)

type BigFile struct {
	ChunkSize         int64            // 分片大小 M
	MaxWorkers        int              // 同时处理最大分块数量，合理用防止超大文件内存益处
	File              string           // 文件路径
	FileMd5           string           // 文件Md5
	ChunkCount        int64            // 分片数量
	HandledChunkCount int64            // 已处理分片数量
	FileChunkCallback func(*FileChunk) // 分片处理消息

	queue      *gomq.BlockingQueue //等待队列
	pool       *gopool.GoPool      //并发池子
	fileReader *os.File
	fileSize   int64
	isFinish   bool

	lock   *golock.GoLock
	__ctx  context.Context
	cancel context.CancelFunc
}

func (b *BigFile) Release() {
	b.pool.StopAndWait()
}

func (b *BigFile) WaitForFinish() {
	<-b.__ctx.Done()
	b.Release()
}

func (b *BigFile) IsFinish() bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.HandledChunkCount == b.ChunkCount
}

func (b *BigFile) Stop() {
	b.cancel()
}

func (b *BigFile) Start() error {
	if b.FileChunkCallback == nil {
		return errors.New("FileChunkCallback is nil")
	}

	e, err := Exist(b.File)
	if err != nil {
		return err
	}

	if !e {
		return errors.New("file not found")
	}

	file, err := os.Open(b.File)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	__ctx, __cancel := context.WithCancel(context.TODO())
	b.__ctx = __ctx
	b.cancel = __cancel
	cSize := float64(fileInfo.Size()) / float64(b.ChunkSize*1024*1024)
	chunkCount := gconv.Int(math.Ceil(cSize))
	b.ChunkCount = int64(chunkCount)
	b.pool = gopool.NewFixedSizePool(b.MaxWorkers, b.MaxWorkers)
	b.queue = gomq.NewBlockingQueue(b.MaxWorkers)
	b.fileSize = fileInfo.Size()
	b.fileReader = file
	b.lock = golock.New()

	goutils.AsyncFunc(func() {
		for i := 0; i < chunkCount; i++ { //把将要处理分片index 压入队列
			b.queue.Enqueue(int64(i))
		}
	})

	//开始处理
	for i := 0; i < b.MaxWorkers; i++ {
		b.NextChunk()
	}

	return nil
}

func (b *BigFile) CheckAllDone() {
	if b.IsFinish() {
		return
	}

	b.lock.WLockFunc(func(parameters ...any) {
		b.HandledChunkCount++
	})

	if b.IsFinish() {
		if b.isFinish == true {
			return
		}
		b.cancel() //全部处理完成
		b.lock.WLockFunc(func(parameters ...any) {
			b.isFinish = true
		})
	}
}

func (b *BigFile) NextChunk() {
	if b.IsFinish() { //全部处理完成
		return
	}

	goutils.AsyncFunc(func() {
		item := b.queue.Dequeue()
		b.pool.Submit(func() {
			chunk, err := b.createChunk(b.fileReader, item.(int64))
			if err != nil {
				return
			}
			b.FileChunkCallback(chunk)
		})
	})
}

func (b *BigFile) createChunk(file *os.File, index int64) (*FileChunk, error) {
	bufferSize := b.ChunkSize * 1024 * 1024 // 每次读取MB
	startPos := index * bufferSize
	buffer := make([]byte, bufferSize)
	fileInfo, _ := file.Stat() //剩下不足一个整个bufferSize ，具体的大小计算出来
	if index == b.ChunkCount-1 {
		buffer = make([]byte, fileInfo.Size()-startPos)
	}

	_, err := file.ReadAt(buffer, startPos) //TODO: 字节读取验证
	if err != nil {
		return nil, err
	}

	hash := goutils.Md5(buffer)
	return &FileChunk{
		Data:             buffer,
		Hash:             hash,
		Index:            index,
		OriginalFileMd5:  b.FileMd5,
		OriginalFileName: fileInfo.Name(),
		FileName:         fmt.Sprintf("%s.part%d", b.FileMd5+filepath.Ext(fileInfo.Name()), index),
	}, nil
}
