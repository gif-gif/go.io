package gozip

import (
	"bytes"
	"compress/gzip"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"

	"io"
	"os"
)

// xNlContentEncoding: br, gzip
const (
	XNlContentEncoding = "X-NL-Content-Encoding" //默认Header 压缩标识
	NOZIP              = "nozip"                 //不压缩
	GZIP               = "gzip"
	BR                 = "br"
	ZSTD               = "zstd"
	GoZipNoType        = "__nozip__" //默认Header 压缩标识
	GoZipType          = "__zip__"   //默认Header 压缩标识
	UnGoZipType        = "__unzip__" //默认Header 压缩标识
)

func GZip(data []byte) ([]byte, error) {
	// 创建一个buffer用于存储压缩后的数据
	var buf bytes.Buffer

	// 创建一个gzip writer
	gzipWriter := gzip.NewWriter(&buf)

	// 写入数据
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}

	// 关闭writer，确保所有数据被刷新到buffer中
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	// 返回压缩后的数据
	return buf.Bytes(), nil
}

func UnGZip(compressedData []byte) ([]byte, error) {
	// 创建一个byte reader
	bytesReader := bytes.NewReader(compressedData)

	// 创建一个gzip reader
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// 读取解压后的数据
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gzipReader); err != nil {
		return nil, err
	}

	// 返回解压后的数据
	return buf.Bytes(), nil
}

func GZipFile(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 创建gzip writer
	gzipWriter := gzip.NewWriter(destFile)
	defer gzipWriter.Close()

	// 从源文件复制到gzip writer
	_, err = io.Copy(gzipWriter, sourceFile)
	return err
}

func UnGZipFile(src, dst string) error {
	// 打开源文件(gzip压缩文件)
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建gzip reader
	gzipReader, err := gzip.NewReader(sourceFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 从gzip reader复制到目标文件
	_, err = io.Copy(destFile, gzipReader)
	return err
}

// BR 压缩
// 压缩数据
func BrZip(data []byte, quality int) ([]byte, error) {
	var buf bytes.Buffer
	// 创建brotli writer，quality参数范围为0-11，值越大压缩率越高但更慢
	brWriter := brotli.NewWriterLevel(&buf, quality)

	// 写入数据
	_, err := brWriter.Write(data)
	if err != nil {
		return nil, err
	}

	// 关闭writer，确保所有数据被写入
	if err := brWriter.Close(); err != nil {
		return nil, err
	}

	// 返回压缩后的数据
	return buf.Bytes(), nil
}

// 解压数据
func UnBrZip(compressedData []byte) ([]byte, error) {
	// 创建brotli reader
	brReader := brotli.NewReader(bytes.NewReader(compressedData))

	// 读取解压后的数据
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, brReader); err != nil {
		return nil, err
	}

	// 返回解压后的数据
	return buf.Bytes(), nil
}

func BrZipFile(src, dst string, quality int) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 创建brotli writer
	brWriter := brotli.NewWriterLevel(destFile, quality)
	defer brWriter.Close()

	// 从源文件复制到brotli writer
	_, err = io.Copy(brWriter, sourceFile)
	return err
}

func UnBrZipFile(src, dst string) error {
	// 打开源文件(brotli压缩文件)
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建brotli reader
	brReader := brotli.NewReader(sourceFile)

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 从brotli reader复制到目标文件
	_, err = io.Copy(destFile, brReader)
	return err
}

func ZstdCompress(inData []byte) ([]byte, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	outData := encoder.EncodeAll(inData, nil)
	return outData, nil
}

func ZstdDecompress(inData []byte) ([]byte, error) {
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	outData, err := decoder.DecodeAll(inData, nil)
	if err != nil {
		return nil, err
	}
	return outData, nil
}
