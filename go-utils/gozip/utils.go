package gozip

import (
	"encoding/binary"
	"github.com/andybalholm/brotli"
	goerror "github.com/gif-gif/go.io/go-error"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/go-utils/gocrypto"
	"time"
)

func Compress(body []byte, compressMethod string, compressType string) (bool, []byte, error) {
	defer goutils.Recovery(func(err any) {
		golog.Warn(err)
	})
	if compressType == "" {
		compressType = UnGoZipType
	}
	if compressMethod == "" { //没有压缩逻辑
		return false, body, nil
	}

	var data []byte
	var err error
	if compressMethod == GZIP { //压缩和解压
		if compressType == GoZipType {
			data, err = GZip(body)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("GZip error ")
			}
		} else {
			data, err = UnGZip(body)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("UnGZip error ")
			}
		}
	} else if compressMethod == BR {
		if compressType == GoZipType { //压缩和解压
			data, err = BrZip(body, brotli.BestCompression)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("BrZip error ")
			}
		} else {
			data, err = UnBrZip(body)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("UnBrZip error ")
			}
		}
	} else { //不 压缩和解压
		data = body
	}

	return true, data, nil
}

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
	defer goutils.Recovery(func(err any) {
		golog.Warn(err)
	})
	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(timestamp))

	randomIv, err := gocrypto.GenerateByteKey(16)
	if err != nil {
		return nil, err
	}

	var compressBytes []byte
	if compressMethod != "" && compressMethod != NOZIP {
		_, compressBytes, err = Compress(data, compressMethod, GoZipType)
		if err != nil {
			return nil, err
		}
	} else {
		compressBytes = data
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
	defer goutils.Recovery(func(err any) {
		golog.Warn(err)
	})
	AesIvLength := 16
	if len(data) < AesIvLength { //非法数据
		return nil, goerror.NewParamErrMsg("非法数据")
	}

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

	if compressMethod != "" && compressMethod != NOZIP {
		_, body, err = Compress(body, compressMethod, UnGoZipType)
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}

// 加密和解密AesCtr(zip(data))
//
// compressMethod 空时不会压缩和解压
func GoDataAesCTRTransformEncode(data []byte, aesKey []byte, aesIv []byte, compressMethod string) ([]byte, error) {
	defer goutils.Recovery(func(err any) {
		golog.Warn(err)
	})

	var compressBytes []byte
	var err error
	if compressMethod != "" && compressMethod != NOZIP {
		_, compressBytes, err = Compress(data, compressMethod, GoZipType)
		if err != nil {
			return nil, err
		}
	} else {
		compressBytes = data
	}

	body, err := gocrypto.AesCTRTransform(compressBytes, aesKey, aesIv)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GoDataAesCTRTransformDecode(data []byte, aesKey []byte, aesIv []byte, compressMethod string) ([]byte, error) {
	defer goutils.Recovery(func(err any) {
		golog.Warn(err)
	})
	var err error
	body, err := gocrypto.AesCTRTransform(data, aesKey, aesIv)
	if err != nil {
		return nil, err
	}

	var compressBytes []byte
	if compressMethod != "" && compressMethod != NOZIP {
		_, compressBytes, err = Compress(body, compressMethod, UnGoZipType)
		if err != nil {
			return nil, err
		}
	} else {
		compressBytes = body
	}
	return compressBytes, nil
}

func GoDataAesCTRTransform(data []byte, aesKey []byte, aesIv []byte, compressMethod string) ([]byte, error) {
	defer goutils.Recovery(func(err any) {
		golog.Warn(err)
	})
	var err error
	body, err := gocrypto.AesCTRTransform(data, aesKey, aesIv)
	if err != nil {
		return nil, err
	}

	var compressBytes []byte
	if compressMethod != "" && compressMethod != NOZIP {
		_, compressBytes, err = Compress(body, compressMethod, compressMethod)
		if err != nil {
			return nil, err
		}
	} else {
		compressBytes = body
	}
	return compressBytes, nil
}
