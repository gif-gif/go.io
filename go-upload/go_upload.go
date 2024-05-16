package goupload

import (
	goutils "github.com/gif-gif/go.io/go-utils"
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
	name, err := goutils.GetFileHeaderMd5Name(file)
	if err != nil {
		return nil, err
	}

	fullName := name + filepath.Ext(file.Filename)
	if ok, _ := goutils.Exist(assetsDir); !ok {
		err := goutils.CreateSavePath(assetsDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	fullName = strings.ToLower(fullName)

	originalPath := filepath.Join(assetsDir, fullName)
	if ok, _ := goutils.Exist(originalPath); !ok {
		err = goutils.SaveFile(file, originalPath)
		if err != nil {
			return nil, err
		}
	}

	res := &FileUploadResult{}
	res.OriginalFile = originalPath
	res.FileName = fullName

	return res, nil
}
