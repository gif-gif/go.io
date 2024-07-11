package gohttpx

import "time"

const (
	TAG = "gohttpx"

	CONTENT_TYPE_XML  = "application/xml"
	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
)

const (
	POST = "post"
	GET  = "get"
)

const (
	HttpUnknownError = 1000
	HttpRetryError   = 2000
	HttpParamsError  = 3000
)

type Request struct {
	Url           string
	Urls          []string // 如果有值，当url 请求失败时继续用这里的接口尝试，直到成功返回或者全部失败
	Method        string
	Body          interface{}       //post body 参数
	QueryParams   map[string]string //get 参数
	Headers       map[string]string
	Timeout       time.Duration
	RetryCount    int
	RetryWaitTime time.Duration

	IsAll         bool //一次性并发，默认false, IsConcurrency=true时生效，isAll=true时，一开始url+urls 并行请求，否则先请求url,再并行请求urls
	IsConcurrency bool //并行处理，默认false，url--> urls 一个一个串行请求
}

type HttpError struct {
	HttpStatusCode int
	Msg            string
	Error          error
	Errors         []*HttpError //重试逻辑的错误列表
}

type res struct{}

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 请求成功
func HttpSuccess(data interface{}) *Response {
	return &Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
}

func HttpSuccessByCode(code int64, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  "ok",
		Data: data,
	}
}

// 请求格式错误，比如参数格式、参数字段名等 不正确
func HttpBadRequest(msg string) *Response {
	return &Response{
		Code: 400,
		Msg:  msg,
		Data: res{},
	}
}

// 用户没有访问权限，需要进行身份认证
func HttpUnauthorized(msg string) *Response {
	return &Response{
		Code: 401,
		Msg:  msg,
		Data: res{},
	}
}

// 用户已进行身份认证，但权限不够
func HttpForbidden(msg string) *Response {
	return &Response{
		Code: 403,
		Msg:  msg,
		Data: res{},
	}
}

// 接口不存在
func HttpNotFound(msg string) *Response {
	return &Response{
		Code: 404,
		Msg:  msg,
		Data: res{},
	}
}

// 服务器内部错误
func HttpServerError(msg string) *Response {
	return &Response{
		Code: 500,
		Msg:  msg,
		Data: res{},
	}
}

// 请求失败
func HttpFail(msg string) *Response {
	return &Response{
		Code: 10001,
		Msg:  msg,
		Data: res{},
	}
}

// 如需返回特殊错误码，调用此接口
func HttpFailForCode(code int64, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: res{},
	}
}
