package goerror

import "fmt"

type HttpCodeError struct {
	errCode  uint32
	errMsg   string
	showType uint32
	traceId  string
	host     string
}

// 返回给前端的错误码
func (e *HttpCodeError) GetErrCode() uint32 {
	return e.errCode
}

// 返回给前端显示端错误信息
func (e *HttpCodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *HttpCodeError) GetShowType() uint32 {
	return e.showType
}

func (e *HttpCodeError) GetTraceId() string {
	return e.traceId
}

func (e *HttpCodeError) GetHost() string {
	return e.host
}

func (e *HttpCodeError) Error() string {
	//return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s，ShowType:%d，TraceId:%s，Host:%s", e.errCode, e.errMsg, e.showType, e.traceId, e.host)
}

type HttpCodeErrorBuilder struct {
	error *HttpCodeError
}

// NewHttpCodeErrorBuilder 创建一个新的 HttpCodeErrorBuilder
func NewHttpCodeErrorBuilder() *HttpCodeErrorBuilder {
	return &HttpCodeErrorBuilder{
		error: &HttpCodeError{},
	}
}

// WithErrCode 设置错误码
func (b *HttpCodeErrorBuilder) WithErrCode(errCode uint32) *HttpCodeErrorBuilder {
	b.error.errCode = errCode
	return b
}

// WithErrMsg 设置错误信息
func (b *HttpCodeErrorBuilder) WithErrMsg(errMsg string) *HttpCodeErrorBuilder {
	b.error.errMsg = errMsg
	return b
}

// WithShowType 设置显示类型
func (b *HttpCodeErrorBuilder) WithShowType(showType uint32) *HttpCodeErrorBuilder {
	b.error.showType = showType
	return b
}

// WithTraceId 设置追踪ID
func (b *HttpCodeErrorBuilder) WithTraceId(traceId string) *HttpCodeErrorBuilder {
	b.error.traceId = traceId
	return b
}

// WithHost 设置主机信息
func (b *HttpCodeErrorBuilder) WithHost(host string) *HttpCodeErrorBuilder {
	b.error.host = host
	return b
}

// Build 构建并返回 CodeError 实例
func (b *HttpCodeErrorBuilder) Build() *HttpCodeError {
	return b.error
}
