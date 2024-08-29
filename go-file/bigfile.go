package gofile

import (
	"context"
	"errors"
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/sync/errgroup"
	"math"
	"os"
	"path/filepath"
)

type BigFile struct {
	ChunkSize           int64                        // 分片大小 M
	MaxWorkers          int                          // 同时处理最大分块数量，合理用防止超大文件内存益处
	File                string                       // 文件路径
	FileMd5             string                       // 文件Md5
	ChunkCount          int64                        // 分片数量
	FileChunkCallback   func(chunk *FileChunk) error // 分片处理消息
	SuccessChunkIndexes []int64                      //处理成功的碎片index

	fileReader   *os.File
	fileSize     int64
	__ctx        context.Context
	errGroupPool *errgroup.Group
}

func (b *BigFile) IsSuccess() bool {
	return len(b.SuccessChunkIndexes) == int(b.ChunkCount)
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

	g, ctx := errgroup.WithContext(context.Background())
	b.__ctx = ctx
	b.errGroupPool = g
	b.errGroupPool.SetLimit(b.MaxWorkers)

	cSize := float64(fileInfo.Size()) / float64(b.ChunkSize*1024*1024)
	chunkCount := gconv.Int(math.Ceil(cSize))
	b.ChunkCount = int64(chunkCount)
	b.fileSize = fileInfo.Size()
	b.fileReader = file

	for i := 0; i < chunkCount; i++ {
		chunkIndex := gconv.Int64(i)
		b.errGroupPool.Go(func() error {
			if goutils.IsContextDone(b.__ctx) {
				return nil
			}
			chunk, err := b.createChunk(b.fileReader, chunkIndex)
			if err != nil {
				return err
			}
			err = b.FileChunkCallback(chunk)
			if err != nil {
				return err
			}
			b.SuccessChunkIndexes = append(b.SuccessChunkIndexes, chunkIndex)
			return nil
		})
	}
	return b.errGroupPool.Wait()
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
