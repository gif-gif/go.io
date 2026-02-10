package gozip

import (
	"github.com/andybalholm/brotli"
	goerror "github.com/gif-gif/go.io/go-error"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
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
	} else if compressMethod == ZSTD {
		if compressType == GoZipType { //压缩和解压
			data, err = ZstdCompress(body)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("ZstdCompress error ")
			}
		} else {
			data, err = ZstdDecompress(body)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("ZstdDecompress error ")
			}
		}
	} else { //不 压缩和解压
		data = body
	}

	return true, data, nil
}
