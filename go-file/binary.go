package gofile

import (
	"bufio"
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
func ReadFileChunks(filePath string, chunkSize int, callback func(chunk []byte) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)
	reader := bufio.NewReader(file)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// 调用回调函数处理数据块
		if err := callback(buffer[:n]); err != nil {
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
