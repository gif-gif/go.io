package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	gofile "github.com/gif-gif/go.io/go-file"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"path/filepath"
)

var uploadPath = "/Users/Jerry/Downloads/chrome/fileparts"
var fileName = "test.apk"

func main() {
	cutFile()
	<-gocontext.Cancel().Done()
}

func cutFile() {
	ts := goutils.MeasureExecutionTime(func() {
		filePath := "/Users/Jerry/Downloads/chrome/dy12.9.0.apk"
		fileMd5, err := goutils.CalculateFileMD5(filePath)
		if err != nil {
			golog.WithTag("gofile").Error(err)
		}

		req := &gofile.BigFile{
			File:       filePath,
			MaxWorkers: 1,
			ChunkSize:  1,
			FileMd5:    fileMd5,
		}

		req.FileChunkCallback = func(chunk *gofile.FileChunk) {
			golog.WithTag("chunkCount").Info(gconv.String(chunk.Index) + ":" + gconv.String(len(chunk.Data)) + ":" + gconv.String(chunk.Hash))
			//存储文件或者上传文件
			chunkFile := filepath.Join(uploadPath, chunk.ChunkFileName)
			err := gofile.WriteToFile(chunkFile, chunk.Data)
			if err != nil {
				req.Stop()
				return
			}

			if !req.IsFinish() { //还有没处理完的继续处理
				req.NextChunk()
			}
			req.DoneOneChunk()
		}

		err = req.Start()
		if err != nil {
			golog.Fatal(err)
		}
		req.WaitForFinish()

		// 调用合并接口
		err = gofile.MergeFile(uploadPath, fileName, fileMd5, req.ChunkCount)
		if err != nil {
			golog.WithTag("gofile").Error(err)
		}
	})

	golog.WithTag("cutFile").Info("执行时间:" + gconv.String(ts))
}
