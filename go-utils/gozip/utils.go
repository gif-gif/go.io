package gozip

import (
	"encoding/binary"
	"github.com/andybalholm/brotli"
	goerror "github.com/gif-gif/go.io/go-error"
	"github.com/gif-gif/go.io/go-utils/gocrypto"
	"time"
)

// 压缩 和 解压逻辑
//func CompressHandler(w http.ResponseWriter, r *http.Request, compressMethod string, compressType string) (bool, []byte, error) {
//	if compressType == "" {
//		compressType = UnGoZipType
//	}
//	if compressMethod == "" { //没有压缩逻辑
//		return false, nil, nil
//	}
//	bytebuffer := bytebufferpool.Get()
//	length, err := bytebuffer.ReadFrom(r.Body)
//	if err != nil {
//		return false, nil, goerror.NewParamErrMsg("ReadFrom error ")
//	}
//	bodyAll := bytebuffer.B[:length]
//	defer bytebufferpool.Put(bytebuffer)
//
//	ok, data, err := Compress(bodyAll, compressMethod, compressType)
//	if err != nil {
//		return false, nil, err
//	}
//	if !ok {
//		return false, data, nil
//	}
//
//	return true, data, nil
//}

func Compress(body []byte, compressMethod string, compressType string) (bool, []byte, error) {
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
	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(timestamp))

	randomIv, err := gocrypto.GenerateByteKey(16)
	if err != nil {
		return nil, err
	}
	_, compressBytes, err := Compress(data, compressMethod, GoZipType)
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

	_, body, err = Compress(body, compressMethod, UnGoZipType)
	if err != nil {
		return nil, err
	}
	
	return body, nil
}
