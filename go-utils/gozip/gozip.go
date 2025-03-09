package gozip

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"github.com/andybalholm/brotli"
	goerror "github.com/gif-gif/go.io/go-error"
	"github.com/gif-gif/go.io/go-utils/gocrypto"
	"io"
	"os"
	"time"
)

// xNlContentEncoding: br, gzip
const (
	XNlContentEncoding = "X-NL-Content-Encoding" //默认Header 压缩标识
	NOZIP              = "nozip"                 //不压缩
	GZIP               = "gzip"
	BR                 = "br"
	GoZipNoType        = "__nozip__" //默认Header 压缩标识
	GoZipType          = "__zip__"   //默认Header 压缩标识
	UnGoZipType        = "__unzip__" //默认Header 压缩标识
)

// AesIv 动态生成 (aes(time+zip(data)))
//
// 1. 先时间戳 time
//
// 2. 生成16位随机IV
//
// 3. 压缩data(如果有压缩) zip(data)
//
// 4. 加密data： aes(time+zip(data))
//
// 5. 拼接iv+encryptData
func GoDataEncrypt(data []byte, AesKey []byte, compressMethod string) ([]byte, error) {
	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(timestamp))

	randomIv, err := gocrypto.GenerateByteKey(16)
	if err != nil {
		return nil, err
	}
	var compressBytes []byte
	if compressMethod == GZIP {
		compressBytes, err = GZip(data)
	} else if compressMethod == BR {
		compressBytes, err = BrZip(data, brotli.BestCompression)
	} else {
		compressBytes = data //不压缩
	}
	if err != nil {
		return nil, err
	}
	timeAndData := append(timestampBytes, compressBytes...)
	dataEncrypt, err := gocrypto.AESCBCEncrypt(timeAndData, AesKey, randomIv)
	if err != nil {
		return nil, err
	}
	resDataEncrypt := append(randomIv, dataEncrypt...)
	return resDataEncrypt, nil
}

// 解密：
//
// 1. 先取前16个字节，作为AES的IV
//
// 2. 取剩余的字节解密
//
// 3. 取前8个字节，作为时间戳
//
// 4. 取剩余的字节
//
// 5. 解压data(如果有压缩)
func GoDataDecrypt(data []byte, AesKey []byte, compressMethod string) ([]byte, error) {
	AesIvLength := 16
	first16BytesIv := data[:AesIvLength]
	// 获取剩余的字节解密
	timeAndDataEncryptBytes := data[AesIvLength:]
	timeAndZipBody, err := gocrypto.AESCBCDecrypt(timeAndDataEncryptBytes, AesKey, first16BytesIv)
	if err != nil {
		return nil, err
	}

	if len(timeAndZipBody) < 8 { //非法数据
		return nil, goerror.NewParamErrMsg("非法数据")
	}

	timestampBytes := timeAndZipBody[:8]
	_ = binary.BigEndian.Uint64(timestampBytes) //客户端时间
	body := timeAndZipBody[8:]

	if compressMethod == GZIP {
		body, err = UnGZip(body)
	} else if compressMethod == BR {
		body, err = UnBrZip(body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

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
