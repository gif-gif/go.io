package gofile

import (
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type FileUploadResult struct {
	OriginalFile string
	FileName     string
}

// 上传file文件到assetsDir目录下，assetsDir 目录不存在则自动创建,返回存储位置
func UploadFile(assetsDir string, file *multipart.FileHeader) (*FileUploadResult, error) {
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

// 上传文件分片
func UploadChunkHandler(assetsDir string, chunkIndex int64, chunkMd5 string, file *multipart.FileHeader) (*FileUploadResult, error) {
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

	if fileMd5 != chunkMd5 { //分片md5 验证
		return nil, fmt.Errorf("fileMd5 mismatch")
	}

	fullName := fileMd5 + filepath.Ext(file.Filename)
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

func MergeFile(filePath string, fileName string, fileMd5 string, totalChunks int) error {
	finalFilePath := filepath.Join(filePath, fileName)
	finalFile, err := os.Create(finalFilePath)
	if err != nil {
		return err
	}
	defer finalFile.Close()

	// 合并所有分片
	for i := 0; i < totalChunks; i++ {
		chunkFilePath := filepath.Join(filePath, fmt.Sprintf("%s.part%d", fileName, i))
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
