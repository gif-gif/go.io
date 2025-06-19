package gocrypto

import (
	"encoding/binary"
	goerror "github.com/gif-gif/go.io/go-error"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/go-utils/gozip"
	"time"
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
	defer goutils.Recovery(func(err any) {
		golog.Warn(err)
	})
	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(timestamp))

	randomIv, err := GenerateByteKey(16)
	if err != nil {
		return nil, err
	}

	var compressBytes []byte
	if compressMethod != "" && compressMethod != gozip.NOZIP {
		_, compressBytes, err = gozip.Compress(data, compressMethod, gozip.GoZipType)
		if err != nil {
			return nil, err
		}
	} else {
		compressBytes = data
	}

	timeAndData := append(timestampBytes, compressBytes...)
	dataEncrypt, err := AESCBCEncrypt(timeAndData, AesKey, randomIv)
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
	timeAndZipBody, err := AESCBCDecrypt(timeAndDataEncryptBytes, AesKey, first16BytesIv)
	if err != nil {
		return nil, err
	}

	if len(timeAndZipBody) < 8 { //非法数据
		return nil, goerror.NewParamErrMsg("非法数据")
	}

	timestampBytes := timeAndZipBody[:8]
	_ = binary.BigEndian.Uint64(timestampBytes) //客户端时间
	body := timeAndZipBody[8:]

	if compressMethod != "" && compressMethod != gozip.NOZIP {
		_, body, err = gozip.Compress(body, compressMethod, gozip.UnGoZipType)
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
	if compressMethod != "" && compressMethod != gozip.NOZIP {
		_, compressBytes, err = gozip.Compress(data, compressMethod, gozip.GoZipType)
		if err != nil {
			return nil, err
		}
	} else {
		compressBytes = data
	}

	body, err := AesCTRTransform(compressBytes, aesKey, aesIv)
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
	body, err := AesCTRTransform(data, aesKey, aesIv)
	if err != nil {
		return nil, err
	}

	var compressBytes []byte
	if compressMethod != "" && compressMethod != gozip.NOZIP {
		_, compressBytes, err = gozip.Compress(body, compressMethod, gozip.UnGoZipType)
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
	body, err := AesCTRTransform(data, aesKey, aesIv)
	if err != nil {
		return nil, err
	}

	var compressBytes []byte
	if compressMethod != "" && compressMethod != gozip.NOZIP {
		_, compressBytes, err = gozip.Compress(body, compressMethod, compressMethod)
		if err != nil {
			return nil, err
		}
	} else {
		compressBytes = body
	}
	return compressBytes, nil
}
