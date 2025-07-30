package goerror

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"runtime"
)

func IsDatabaseNoRowsError(err error) bool {
	return err != nil && err.Error() == "sql: no rows in result set"
}

func GetStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func IsErrCode(err error, code uint32) bool {
	errCode, _ := GetErrCodeMsg(err)
	return code == errCode
}

// 建议用 GetCodeError
func GetErrCodeMsg(err error) (errCode uint32, errMsg string) {
	if err == nil {
		return 0, ""
	}
	errCode = SERVER_COMMON_ERROR
	errMsg = "server error"
	causeErr := errors.Cause(err)           // err类型
	if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
		errCode = e.GetErrCode()
		errMsg = e.GetErrMsg()
	} else if er, ok := err.(*HttpCodeError); ok { //http 状态码
		errCode = er.GetErrCode()
		errMsg = er.GetErrMsg()
	} else { //通用错误
		errCode, errMsg = GetErrorMsg(err)
	}
	return
}

func GetErrorMsg(err error) (uint32, string) {
	causeErr := errors.Cause(err)
	if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
		grpcCode := uint32(gstatus.Code())
		return grpcCode, gstatus.Message()
	}
	return 500, "server error"
}

// CodeErrorBuilder.build() 构建 CodeError
//
// 返回错误码 CodeErrorBuilder
func GetCodeError(err error) *CodeErrorBuilder {
	codeErr := NewCodeErrorBuilder()
	if err == nil {
		return codeErr
	}
	errCode, errMsg := GetErrCodeMsg(err)
	codeErr.WithErrCode(errCode).WithErrMsg(errMsg)
	return codeErr
}
