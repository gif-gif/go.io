package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	gofile "github.com/gif-gif/go.io/go-file"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"time"
)

func main() {
	test()
	<-gocontext.Cancel().Done()
}

func test() {
	ts := goutils.MeasureExecutionTime(func() {
		req := &gofile.BigFileRequest{
			File:       "/Users/Jerry/Downloads/chrome/dy12.9.0.apk",
			MaxWorkers: 10,
			ChunkSize:  10,
		}

		req.FileChunkCallback = func(chunk *gofile.FileChunk) {
			golog.WithTag("chunkCount").Info(gconv.String(chunk.Index) + ":" + gconv.String(chunk.ByteLength) + ":" + gconv.String(chunk.Hash))
			time.Sleep(1 * time.Second) //模拟分片处理耗时
			if !req.IsFinish() {        //还有没处理完的继续处理
				req.NextChunk()
			}
			req.DoneOneChunk()
		}

		err := req.Start()
		if err != nil {
			golog.Fatal(err)
		}
		req.WaitForFinish()
		req.Release()
	})

	golog.WithTag("file").Info("执行时间:" + gconv.String(ts))
}
