package gofile

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func UnzipFile(zipFile, destDir string) error {
	// 打开 zip 文件
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// 遍历 zip 文件内容
	for _, file := range reader.File {
		// 构建完整路径
		path := filepath.Join(destDir, file.Name)

		// 如果是目录，创建它
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		// 创建目标文件
		destFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		// 打开 zip 中的文件
		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return err
		}

		// 复制内容
		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
