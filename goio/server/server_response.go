package goserver

import (
	"encoding/json"
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/go-playground/validator/v10"

	"strings"
)

type Response struct {
	Success      bool   `json:"success"`
	Data         any    `json:"data"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ShowType     uint32 `json:"showType"`
	TraceId      string `json:"traceId"`
	Host         string `json:"host"`
}

// ResponseBuilder 是 Response 的构建器
type ResponseBuilder struct {
	response *Response
}

// NewResponseBuilder 创建一个新的 Response 构建器
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		response: &Response{
			Success: true, // 默认设置为成功
		},
	}
}

// WithSuccess 设置成功状态
func (b *ResponseBuilder) WithSuccess(success bool) *ResponseBuilder {
	b.response.Success = success
	return b
}

// WithData 设置数据
func (b *ResponseBuilder) WithData(data any) *ResponseBuilder {
	b.response.Data = data
	return b
}

// WithErrorCode 设置错误代码
func (b *ResponseBuilder) WithErrorCode(errorCode string) *ResponseBuilder {
	b.response.ErrorCode = errorCode
	return b
}

// WithErrorMessage 设置错误信息
func (b *ResponseBuilder) WithErrorMessage(errorMessage string) *ResponseBuilder {
	b.response.ErrorMessage = errorMessage
	return b
}

// WithShowType 设置显示类型
func (b *ResponseBuilder) WithShowType(showType uint32) *ResponseBuilder {
	b.response.ShowType = showType
	return b
}

// WithTraceId 设置追踪ID
func (b *ResponseBuilder) WithTraceId(traceId string) *ResponseBuilder {
	b.response.TraceId = traceId
	return b
}

// WithHost 设置主机信息
func (b *ResponseBuilder) WithHost(host string) *ResponseBuilder {
	b.response.Host = host
	return b
}

// Build 构建并返回最终的 Response 对象
func (b *ResponseBuilder) Build() *Response {
	return b.response
}

// 便捷方法: 快速创建成功响应
func SuccessResponse(data any) *Response {
	return NewResponseBuilder().
		WithSuccess(true).
		WithData(data).
		Build()
}

// 便捷方法: 快速创建错误响应
func ErrorResponse(errorCode string, errorMessage string, showType uint32) *Response {
	return NewResponseBuilder().
		WithSuccess(false).
		WithErrorCode(errorCode).
		WithErrorMessage(errorMessage).
		WithShowType(showType).
		Build()
}

func (rsp *Response) Copy() *Response {
	r := NewResponseBuilder().
		WithSuccess(rsp.Success).
		WithData(rsp.Data).
		WithErrorCode(rsp.ErrorCode).
		WithErrorMessage(rsp.ErrorMessage).
		WithShowType(rsp.ShowType).
		WithTraceId(rsp.TraceId).
		WithHost(rsp.Host).
		Build()
	return r
}

func (rsp *Response) String() string {
	buf, err := json.Marshal(rsp)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

// Error 创建带格式化错误消息的错误响应
func Error(code int32, message string, v ...interface{}) *Response {
	return NewResponseBuilder().
		WithSuccess(false).
		WithErrorCode(fmt.Sprintf("%d", code)).
		WithErrorMessage(fmt.Sprintf(message, v...)).
		Build()
}

func ErrorWithValidate(err error, messages map[string]string) *Response {
	if v, ok := err.(*json.UnmarshalTypeError); ok {
		return Error(7001, fmt.Sprintf("请求参数 %s 的类型是 %s, 不是 %s", v.Field, v.Type, v.Value))
	}

	if v, ok := err.(validator.ValidationErrors); ok {
		for _, i := range v {
			field := goutils.Camel2Case(i.Field())
			key := fmt.Sprintf("%s_%s", field, strings.ToLower(i.Tag()))
			if msg, ok := messages[key]; ok {
				return Error(7002, msg)
			}
			return Error(7003, fmt.Sprintf("%s %s", field, i.Tag()))
		}
	}

	return Error(7004, "参数错误", err)
}
