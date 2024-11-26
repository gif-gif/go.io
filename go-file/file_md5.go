package gofile

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"io"
	"io/ioutil"
	"os"
)

// 计算文件MD5 支持大文件
func MD5(file string) (string, error) {
	defaultSize := int64(16 * 1024 * 1024)

	info, err := os.Stat(file)
	if err != nil {
		golog.Error(err)
		return "", err
	}

	// 小文件
	if info.Size() < defaultSize {
		b, _ := ioutil.ReadFile(file)
		sum := md5.Sum(b)
		return hex.EncodeToString(sum[:]), nil
	}

	// 大文件
	{
		tempFile, err := ioutil.TempFile(os.TempDir(), "goo-md5-temp-file")
		if err != nil {
			golog.Error(err)
			return "", err
		}
		defer tempFile.Close()

		f, err := os.OpenFile(file, os.O_RDONLY, 0755)
		if err != nil {
			golog.Error(err)
			return "", err
		}
		defer f.Close()

		io.Copy(tempFile, f)
		tempFile.Seek(0, os.SEEK_SET)

		h := md5.New()
		io.Copy(h, tempFile)
		return hex.EncodeToString(h.Sum(nil)), nil
	}
}

// 计算文件md5(支持超大文件)
func CalculateFileMD5(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建MD5哈希对象
	hash := md5.New()

	// 创建一个缓冲区，逐块读取文件内容
	buffer := make([]byte, 1024*1024) // 1MB 缓冲区
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		// 更新哈希值
		if _, err := hash.Write(buffer[:n]); err != nil {
			return "", err
		}
	}

	// 计算最终的哈希值
	hashInBytes := hash.Sum(nil)
	hashInString := fmt.Sprintf("%x", hashInBytes)

	return hashInString, nil
}
