package goerror

import (
	"fmt"
)

// 常用通用固定错误
//
//	codeError := NewCodeErrorBuilder().
//		WithErrCode(1001).
//		WithErrMsg("系统错误").
//		WithShowType(1).
//		WithTraceId("trace-123").
//		WithHost("localhost").
//		Build()
type CodeError struct {
	errCode  uint32
	errMsg   string
	showType uint32
	traceId  string
	host     string
}

// CodeErrorBuilder 是 CodeError 的构建器
type CodeErrorBuilder struct {
	error *CodeError
}

// 返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) GetShowType() uint32 {
	return e.showType
}

func (e *CodeError) GetTraceId() string {
	return e.traceId
}

func (e *CodeError) GetHost() string {
	return e.host
}

func (e *CodeError) Error() string {
	//return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s，ShowType:%d，TraceId:%s，Host:%s", e.errCode, e.errMsg, e.showType, e.traceId, e.host)
}

func NewErrCodeMsg(errCode uint32, errMsg string) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(errMsg)
	return builder.Build()
}

func NewErrCodeMsgForMessageError(errCode uint32, errMsg string) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(errMsg).WithShowType(ShowTypeMessageError)
	return builder.Build()
}

func NewErrCodeMsgForNotification(errCode uint32, errMsg string) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(errMsg).WithShowType(ShowTypeNotification)
	return builder.Build()
}

func NewErrCodeMsgForMessageWarn(errCode uint32, errMsg string) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(errMsg).WithShowType(ShowTypeMessageWarn)
	return builder.Build()
}

func NewErrCodeMsgForPage(errCode uint32, errMsg string) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(errMsg).WithShowType(ShowTypePage)
	return builder.Build()
}

func NewErrCode(errCode uint32) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(MapErrMsg(errCode))
	return builder.Build()
}

func NewErrMsg(errMsg string) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(SERVER_COMMON_ERROR).WithErrMsg(errMsg)
	return builder.Build()
}

func NewParamErrMsg(errMsg string) *CodeError {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(REQUEST_PARAM_ERROR).WithErrMsg(errMsg)
	return builder.Build()
}

func NewErrorMsg(errCode uint32, errMsg string) error {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(errMsg)
	return builder.Build()
}

func NewError(errCode uint32) error {
	builder := NewCodeErrorBuilder()
	builder.WithErrCode(errCode).WithErrMsg(MapErrMsg(errCode))
	return builder.Build()
}

// NewCodeErrorBuilder 创建一个新的 CodeErrorBuilder
func NewCodeErrorBuilder() *CodeErrorBuilder {
	return &CodeErrorBuilder{
		error: &CodeError{},
	}
}

// WithErrCode 设置错误码
func (b *CodeErrorBuilder) WithErrCode(errCode uint32) *CodeErrorBuilder {
	b.error.errCode = errCode
	return b
}

// WithErrMsg 设置错误信息
func (b *CodeErrorBuilder) WithErrMsg(errMsg string) *CodeErrorBuilder {
	b.error.errMsg = errMsg
	return b
}

// WithShowType 设置显示类型
func (b *CodeErrorBuilder) WithShowType(showType uint32) *CodeErrorBuilder {
	b.error.showType = showType
	return b
}

// WithTraceId 设置追踪ID
func (b *CodeErrorBuilder) WithTraceId(traceId string) *CodeErrorBuilder {
	b.error.traceId = traceId
	return b
}

// WithHost 设置主机信息
func (b *CodeErrorBuilder) WithHost(host string) *CodeErrorBuilder {
	b.error.host = host
	return b
}

// Build 构建并返回 CodeError 实例
func (b *CodeErrorBuilder) Build() *CodeError {
	return b.error
}
