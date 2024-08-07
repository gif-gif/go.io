package gofile

import (
	"errors"
	"fmt"
	gohttpx "github.com/gif-gif/go.io/go-http/go-httpex"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gogf/gf/util/gconv"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const (
	UPLOAD_TYPE_LOCAL = 1
	UPLOAD_TYPE_OSS   = 2
	UPLOAD_TYPE_CHUNK = 3
)

type FileUploadResult struct {
	OriginalFile string
	FileName     string
}

// 服务器接受file文件到assetsDir目录下，assetsDir 目录不存在则自动创建,返回存储位置
func ReceiveFile(assetsDir string, file *multipart.FileHeader) (*FileUploadResult, error) {
	if ok, _ := Exist(assetsDir); !ok {
		err := CreateSavePath(assetsDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	fileMd5, err := GetFileHeaderMd5Name(file)
	if err != nil {
		return nil, err
	}

	fullName := fileMd5 + filepath.Ext(file.Filename)
	fullName = strings.ToLower(fullName)

	originalPath := filepath.Join(assetsDir, fullName)
	if ok, _ := Exist(originalPath); !ok {
		err = SaveFile(file, originalPath)
		if err != nil {
			return nil, err
		}
	}

	res := &FileUploadResult{}
	res.OriginalFile = originalPath
	res.FileName = fullName

	return res, nil
}

// 服务器接受文件分片
func ReceiveChunkHandler(assetsDir string, chunkIndex int64, chunkMd5 string, fileMd5 string, file *multipart.FileHeader) (*FileUploadResult, error) {
	if ok, _ := Exist(assetsDir); !ok {
		err := CreateSavePath(assetsDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	chunkFileMd5, err := GetFileHeaderMd5Name(file)
	if err != nil {
		return nil, err
	}

	if chunkFileMd5 != chunkMd5 { //分片md5 验证
		return nil, fmt.Errorf("chunkFileMd5 mismatch")
	}

	fullName := fileMd5 + filepath.Ext(file.Filename) //以原文件md5作为命名
	fullName = strings.ToLower(fullName)
	// 创建分片文件
	chunkFilePath := filepath.Join(assetsDir, fmt.Sprintf("%s.part%d", fullName, chunkIndex))

	if ok, _ := Exist(chunkFilePath); !ok {
		err = SaveFile(file, chunkFilePath)
		if err != nil {
			return nil, err
		}
	}

	res := &FileUploadResult{}
	res.OriginalFile = chunkFilePath
	res.FileName = fullName
	return res, nil
}

// 服务器合并所有文件分片，并验证md5
func MergeFile(filePath string, fileName string, fileMd5 string, totalChunks int64) error {
	fileName = fileMd5 + filepath.Ext(fileName)
	finalFilePath := filepath.Join(filePath, fileName)
	finalFile, err := os.Create(finalFilePath)
	if err != nil {
		return err
	}
	defer finalFile.Close()

	// 合并所有分片
	for i := 0; i < int(totalChunks); i++ {
		chunkFilePath := filepath.Join(filePath, fmt.Sprintf("%s.part%d", fileMd5+filepath.Ext(fileName), i))
		chunkFile, err := os.Open(chunkFilePath)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		_, err = io.Copy(finalFile, chunkFile)
		if err != nil {
			return err
		}

		// 删除分片文件
		err = os.Remove(chunkFilePath)
		if err != nil {
			return err
		}
	}

	finalFileMd5, err := goutils.CalculateFileMD5(finalFilePath)
	if err != nil {
		return err
	}

	if fileMd5 != finalFileMd5 { //最终文件md5验证
		return fmt.Errorf("fileMd5 mismatch")
	}

	//golog.WithTag("mergeFile").Info("执行时间:" + gconv.String(ts))
	return nil
}

// 上传一个文件分片
func UploadChunk(url string, fileName string, fileMd5 string, chunkIndex int, chunkData []byte) (*gohttpx.Response, error) {
	req := &gohttpx.Request{
		Url:    url,
		Method: gohttpx.POST,
		Body:   chunkData,
		FormData: map[string]string{
			"type":       gconv.String(UPLOAD_TYPE_CHUNK),
			"fileName":   fileName,
			"fileMd5":    fileMd5,
			"chunkIndex": gconv.String(chunkIndex),
		},
	}

	res := &gohttpx.Response{}
	err := gohttpx.HttpPost[gohttpx.Response](req, res)
	if err != nil {
		return nil, errors.New(err.ErrorInfo())
	}
	return res, nil
}
