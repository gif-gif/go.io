package gofile

import (
	"context"
	"errors"
	golock "github.com/gif-gif/go.io/go-lock"
	golog "github.com/gif-gif/go.io/go-log"
	gomq "github.com/gif-gif/go.io/go-mq"
	gopool "github.com/gif-gif/go.io/go-pool"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"math"
	"os"
)

// fmt.Sprintf("%s.part%d", fileName, i)
// mapreduce big file
// 大文件逻辑 for 把大文件并发分片处理，为了防止OOM超大文件边分片边处理的策略
type FileChunk struct {
	Data       []byte //分片数据
	Hash       string //分片Hash
	Index      int64  //分片顺序号
	ByteLength int64  //分片大小 len(Data)
}

type BigFile struct {
	ChunkSize         int64            // 分片大小 M
	MaxWorkers        int              // 同时处理最大分块数量，合理用防止超大文件内存益处
	File              string           // 文件路径
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

func (b *BigFile) DoneOneChunk() {
	if b.IsFinish() {
		golog.WithTag("DoneOneChunk").Info("Chunk count is already finished")
		return
	}

	b.lock.WLockFunc(func(parameters ...any) {
		b.HandledChunkCount++
	})

	if b.IsFinish() {
		if b.isFinish == true {
			return
		}
		golog.WithTag("DoneOneChunk").Info("cancel")
		b.cancel() //全部处理完成
		b.lock.WLockFunc(func(parameters ...any) {
			b.isFinish = true
		})
	}
}

func (b *BigFile) NextChunk() {
	if b.IsFinish() { //全部处理完成
		golog.WithTag("NextChunk").Info("Chunk count is already finished")
		return
	}

	goutils.AsyncFunc(func() {
		item := b.queue.Dequeue()
		b.pool.Submit(func() {
			chunk, err := createChunk(b.fileReader, item.(int64), b.ChunkSize, b.ChunkCount)
			if err != nil {
				return
			}
			b.FileChunkCallback(chunk)
		})
	})
}

func createChunk(file *os.File, index int64, chunkSize int64, chunkCount int64) (*FileChunk, error) {
	bufferSize := chunkSize * 1024 * 1024 // 每次读取MB
	startPos := index * bufferSize
	buffer := make([]byte, bufferSize)
	if index == chunkCount-1 { //剩下不足一个整个bufferSize ，具体的大小计算出来
		fileInfo, _ := file.Stat()
		buffer = make([]byte, fileInfo.Size()-startPos)
	}

	n, err := file.ReadAt(buffer, startPos) //TODO: 字节读取验证
	if err != nil {
		return nil, err
	}

	return &FileChunk{
		Data:       buffer,
		Hash:       goutils.Md5(buffer),
		Index:      index,
		ByteLength: int64(n),
	}, nil
}
