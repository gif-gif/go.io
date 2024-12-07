package gofile

import (
	"io"
	"os"
)

// ReadEntireFile 读取整个文件到字节切片
func ReadEntireFile(filePath string) ([]byte, error) {
	// 使用 os.ReadFile，这是最简单的方法
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReadFileChunks 分块读取大文件，避免内存占用过大
func ReadFileChunks(filename string, chunkSize int, callback func(chunk []byte) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用bytes.Buffer来读取文件
	for {
		// 每次创建新的buffer
		chunk := make([]byte, chunkSize)
		_, err := file.Read(chunk)
		// 调用回调函数处理数据块
		// 只传递实际读取的数据
		if err := callback(chunk); err != nil {
			return err
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// ReadFileAt 从指定位置读取文件
func ReadFileAt(filePath string, offset int64, length int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, length)
	_, err = file.ReadAt(data, offset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return data, nil
}
