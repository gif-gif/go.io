package goerror

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"runtime"
)

func GetStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func IsErrCode(err error, code uint32) bool {
	errCode, _ := GetErrCodeMsg(err)
	return code == errCode
}

func GetErrCodeMsg(err error) (errCode uint32, errMsg string) {
	if err == nil {
		return 0, ""
	}
	errCode = 500
	errMsg = "server error"

	causeErr := errors.Cause(err)           // err类型
	if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
		//自定义CodeError
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
	} else {
		if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
			grpcCode := uint32(gstatus.Code())
			if IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
				errCode = grpcCode
				errMsg = gstatus.Message()
			} else {
				if errorsx != nil {
					if _, ok := errorsx[grpcCode]; ok {
						errCode = grpcCode
						errMsg = gstatus.Message()
					}
				}
			}

		}
	}

	return
}
