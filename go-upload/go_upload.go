package goupload

import (
	goutils "github.com/jiriyao/go.io/go-utils"
	"mime/multipart"
	"os"
	"path/filepath"
)

// 上传file文件到assetsDir目录下，assetsDir 目录不存在则自动创建
func UploadFile(assetsDir string, file *multipart.FileHeader) error {
	name, err := goutils.GetFileHeaderMd5Name(file)

	if err != nil {
		return err
	}

	fullName := name + filepath.Ext(file.Filename)
	if ok, _ := goutils.Exist(assetsDir); !ok {
		err := goutils.CreateSavePath(assetsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	originalPath := filepath.Join(assetsDir, fullName)
	if ok, _ := goutils.Exist(originalPath); !ok {
		err = goutils.SaveFile(file, originalPath)
		if err != nil {
			return err
		}
	}

	return nil
}
