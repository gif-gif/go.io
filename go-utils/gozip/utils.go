package gozip

import (
	"github.com/andybalholm/brotli"
	goerror "github.com/gif-gif/go.io/go-error"
	"github.com/valyala/bytebufferpool"
	"net/http"
)

// 压缩 和 解压逻辑
func CompressHandler(w http.ResponseWriter, r *http.Request, compressMethod string, compressType string) (bool, []byte, error) {
	if compressType == "" {
		compressType = UnGoZipType
	}
	if compressMethod == "" { //没有压缩逻辑
		return false, nil, nil
	}
	bytebuffer := bytebufferpool.Get()
	length, err := bytebuffer.ReadFrom(r.Body)
	if err != nil {
		return false, nil, goerror.NewParamErrMsg("ReadFrom error ")
	}
	bodyAll := bytebuffer.B[:length]
	defer bytebufferpool.Put(bytebuffer)

	var data []byte
	if compressMethod == GZIP { //压缩和解压
		if compressType == GoZipType {
			data, err = GZip(bodyAll)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("GZip error ")
			}
		} else {
			data, err = UnGZip(bodyAll)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("UnGZip error ")
			}
		}
	} else {
		if compressType == GoZipType { //压缩和解压
			data, err = BrZip(bodyAll, brotli.BestCompression)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("BrZip error ")
			}
		} else {
			data, err = UnBrZip(bodyAll)
			if err != nil {
				return false, nil, goerror.NewParamErrMsg("UnBrZip error ")
			}
		}
	}

	return true, data, nil
}
