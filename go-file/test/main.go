package main

import (
	"errors"
	gocontext "github.com/gif-gif/go.io/go-context"
	gofile "github.com/gif-gif/go.io/go-file"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
)

var uploadPath = "/Users/Jerry/Downloads/chrome/fileparts"
var fileName = "test.apk"

func main() {
	cutLocalFile()
	<-gocontext.WithCancel().Done()
}

func cutHttpFile() {
	ts := goutils.MeasureExecutionTime(func() {
		filePath := "/Users/Jerry/Downloads/chrome/dy12.9.0.apk"
		fileMd5, err := goutils.CalculateFileMD5(filePath)
		if err != nil {
			golog.WithTag("gofile").Error(err)
		}

		req := &gofile.BigFile{
			File:       filePath,
			MaxWorkers: 3,
			ChunkSize:  1,
			FileMd5:    fileMd5,
		}

		req.FileChunkCallback = func(chunk *gofile.FileChunk) error {
			golog.WithTag("chunkCount").Info(gconv.String(chunk.Index) + ":" + gconv.String(len(chunk.Data)) + ":" + gconv.String(chunk.Hash))
			//存储文件或者上传文件
			rst, err := gofile.UploadChunk("http://localhost:20085/bot/api/file-uploader", chunk)
			if err != nil {
				return err
			}

			if rst.Code != 0 {
				golog.WithTag("gofile").Error(rst.Code)
				return errors.New(gconv.String(rst.Code) + ":" + rst.Msg)
			}

			//time.Sleep(1 * time.Second) // for test
			return nil
		}

		err = req.Start()
		if err != nil {
			golog.WithTag("gofile").Fatal(err)
		}

		if req.IsSuccess() {
			golog.WithTag("gofile").Info("已处理分片:", len(req.SuccessChunkIndexes))
		} else {
			golog.WithTag("gofile").Error("分片上传失败")
			return
		}

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

	golog.WithTag("cutHttpFile").Info("执行时间:" + gconv.String(ts))
}

func cutLocalFile() {
	ts := goutils.MeasureExecutionTime(func() {
		filePath := "/Users/Jerry/Downloads/chrome/dy12.9.0.apk"
		fileMd5, err := goutils.CalculateFileMD5(filePath)
		if err != nil {
			golog.WithTag("gofile").Error(err)
		}

		req := &gofile.BigFile{
			File:       filePath,
			MaxWorkers: 3,
			ChunkSize:  1,
			FileMd5:    fileMd5,
		}

		req.FileChunkCallback = func(chunk *gofile.FileChunk) error {
			golog.WithTag("chunkCount").Info(gconv.String(chunk.Index) + ":" + gconv.String(len(chunk.Data)) + ":" + gconv.String(chunk.Hash))
			//存储文件或者上传文件
			return gofile.SaveToLocal(uploadPath, chunk)
		}

		err = req.Start()
		if err != nil {
			golog.WithTag("gofile").Fatal(err)
		}

		if req.IsSuccess() {
			golog.WithTag("gofile").Info("已处理分片:", len(req.SuccessChunkIndexes))
		} else {
			golog.WithTag("gofile").Error("分片上传失败")
			return
		}

		// 调用合并接口
		rst, err := gofile.MergeFileForChunks(uploadPath, fileName, fileMd5, req.ChunkCount, false)
		if err != nil {
			golog.WithTag("gofile").Error(err.Error())
			return
		}

		if err != nil {
			golog.WithTag("gofile").Error(err)
		}

		golog.WithTag("gofile").Info(rst)
	})
	golog.WithTag("cutHttpFile").Info("执行时间:" + gconv.String(ts))
}
