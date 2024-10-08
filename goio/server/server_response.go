package goserver

import (
	"encoding/json"
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/go-playground/validator/v10"

	"strings"
)

type Response struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Errors  []interface{} `json:"-"`
}

func (rsp *Response) Copy() *Response {
	r := &Response{
		Code:    rsp.Code,
		Message: rsp.Message,
		Data:    rsp.Data,
		Errors:  rsp.Errors,
	}
	return r
}

func (rsp *Response) String() string {
	buf, err := json.Marshal(rsp)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func Success(data interface{}) *Response {
	if data == nil {
		data = map[string]interface{}{}
	}
	return &Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

func Error(code int32, message string, v ...interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Errors:  v,
	}
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
