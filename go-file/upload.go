package gofile

import (
	"context"
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

// 服务器接受file文件到assetsDir目录下，assetsDir 目录不存在则自动创建,返回存储位置
func ReceiveFile(assetsDir string, file *multipart.FileHeader) (*FileReceiveResult, error) {
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

	res := &FileReceiveResult{}
	res.OriginalFile = originalPath
	res.FileName = fullName

	return res, nil
}

// 服务器接受文件分片
func ReceiveChunkHandler(assetsDir string, chunkIndex int64, chunkMd5 string, fileMd5 string, file *multipart.FileHeader) (*FileReceiveResult, error) {
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

	res := &FileReceiveResult{}
	res.OriginalFile = chunkFilePath
	res.FileName = fullName
	return res, nil
}

// 服务器合并所有文件分片，并验证md5, isNotRemoveChunk =true 合并后时不会删除分片
func MergeFileForChunks(filePath string, fileName string, fileMd5 string, totalChunks int64, isNotRemoveChunk bool) (*FileReceiveResult, error) {
	fileName = fileMd5 + filepath.Ext(fileName)
	finalFilePath := filepath.Join(filePath, fileName)
	finalFile, err := os.Create(finalFilePath)
	if err != nil {
		return nil, err
	}
	defer finalFile.Close()

	// 合并所有分片
	for i := 0; i < int(totalChunks); i++ {
		chunkFilePath := filepath.Join(filePath, fmt.Sprintf("%s.part%d", fileMd5+filepath.Ext(fileName), i))
		chunkFile, err := os.Open(chunkFilePath)
		if err != nil {
			return nil, err
		}
		defer chunkFile.Close()

		_, err = io.Copy(finalFile, chunkFile)
		if err != nil {
			return nil, err
		}

		if !isNotRemoveChunk {
			// 删除分片文件
			err = os.Remove(chunkFilePath)
			if err != nil {
				return nil, err
			}
		}
	}

	finalFileMd5, err := goutils.CalculateFileMD5(finalFilePath)
	if err != nil {
		return nil, err
	}

	if fileMd5 != finalFileMd5 { //最终文件md5验证
		return nil, fmt.Errorf("fileMd5 mismatch")
	}

	res := &FileReceiveResult{
		ChunkCount:   totalChunks,
		OriginalFile: finalFilePath,
		FileName:     fileName,
	}

	//golog.WithTag("mergeFile").Info("执行时间:" + gconv.String(ts))
	return res, nil
}

// ////////////////////////////////////////////////////////////// http server upload and merge 供参考
// 上传一个文件分片，（作为客户端请求时验证非法请求认证逻辑需要加，如authToken sign 等等）
func UploadChunk(url string, chunk *FileChunk) (*gohttpx.Response, error) {
	req := &gohttpx.Request{
		Url:       url,
		Method:    gohttpx.POST,
		FileBytes: chunk.Data,
		MultipartFormData: map[string]string{
			"type":       gconv.String(UPLOAD_TYPE_CHUNK),
			"fileName":   chunk.FileName,
			"fileMd5":    chunk.OriginalFileMd5,
			"chunkMd5":   chunk.Hash,
			"chunkIndex": gconv.String(chunk.Index),
		},
		FileName: chunk.OriginalFileName,
		Headers:  map[string]string{"User-Agent": "github.com/gif-gif/go.io"},
	}

	res := &gohttpx.Response{}
	err := gohttpx.HttpPost[gohttpx.Response](context.Background(), req, res)
	if err != nil {
		return nil, errors.New(err.ErrorInfo())
	}
	return res, nil
}

// 分片全部上传完毕后，再请求文件分片合并请求（作为客户端请求时验证非法请求认证逻辑需要加，如authToken sign 等等）
func MergeChunk(url string, fileMergeReq *FileMergeReq) (*gohttpx.Response, error) {
	req := &gohttpx.Request{
		Url:     url,
		Method:  gohttpx.POST,
		Headers: map[string]string{"User-Agent": "github.com/gif-gif/go.io"},
		Body:    fileMergeReq,
	}

	res := &gohttpx.Response{}
	err := gohttpx.HttpPost[gohttpx.Response](context.Background(), req, res)
	if err != nil {
		return nil, errors.New(err.ErrorInfo())
	}
	return res, nil
}

// //////////////////////////////////////////////////////////////local save and merge
// 把分片存在指定目录
func SaveToLocal(savePath string, chunk *FileChunk) error {
	chunkFile := filepath.Join(savePath, chunk.FileName)
	err := WriteToFile(chunkFile, chunk.Data)
	if err != nil {
		return err
	}
	return nil
}
