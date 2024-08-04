package gofile

import (
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
	name, err := GetFileHeaderMd5Name(file)
	if err != nil {
		return nil, err
	}

	fullName := name + filepath.Ext(file.Filename)
	if ok, _ := Exist(assetsDir); !ok {
		err := CreateSavePath(assetsDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

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

//TODO: 大文件上传逻辑
