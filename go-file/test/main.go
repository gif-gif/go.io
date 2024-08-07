package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	gofile "github.com/gif-gif/go.io/go-file"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
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
			rst, err := gofile.UploadChunk("http://localhost:20085/bot/api/file-uploader", chunk)
			if err != nil {
				req.Stop()
				return
			}

			if rst.Code != 0 {
				req.Stop()
				golog.WithTag("gofile").Error(rst.Code)
				return
			}
			if !req.IsFinish() { //还有没处理完的继续处理
				req.NextChunk()
			}
			req.CheckAllDone()
		}

		err = req.Start()
		if err != nil {
			golog.Fatal(err)
		}

		req.WaitForFinish()

		rst, err := gofile.MergeChunk("http://localhost:20085/bot/api/file-merge-uploader", &gofile.FileMergeReq{
			FileMd5:     fileMd5,
			TotalChunks: req.ChunkCount,
			FileName:    fileName,
		})

		if err != nil {
			golog.WithTag("gofile").Error(err.Error())
			return
		}

		if rst.Code != 0 {
			golog.WithTag("gofile").Error(rst.Code)
			return
		}

		// 调用合并接口
		//_, err = gofile.MergeFileForChunks(uploadPath, fileName, fileMd5, req.ChunkCount)
		if err != nil {
			golog.WithTag("gofile").Error(err)
		}
	})

	golog.WithTag("cutFile").Info("执行时间:" + gconv.String(ts))
}
