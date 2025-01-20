package gohttp

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/go-resty/resty/v2"
	"net/url"
	"time"
)

const (
	TAG = "gohttp"

	CONTENT_TYPE_XML  = "application/xml"
	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
)

const (
	POST   = "post"
	GET    = "get"
	PUT    = "put"
	DELETE = "delete"
)

const (
	HttpUnknownError = 1000
	HttpRetryError   = 2000
	HttpParamsError  = 3000
)

type Request struct {
	Url          string
	Urls         []string // 如果有值，当url 请求失败时继续用这里的接口尝试，直到成功返回或者全部失败
	Method       string
	Body         interface{}       //post body 参数
	QueryParams  map[string]string //get 参数
	ParamsValues url.Values        //get 参数
	FormData     map[string]string //formdata 参数
	Headers      map[string]string

	Files             map[string]string //上传文件列表
	FileName          string            //文件名称
	MultipartFormData map[string]string
	FileBytes         []byte

	Timeout       time.Duration
	RetryCount    int
	RetryWaitTime time.Duration
	proxyURL      string

	SetCloseConnection bool //是否关闭连接

	//gotrace infos
	TraceInfo     resty.TraceInfo
	ResponseTime  time.Duration
	ResponseProto string
	Response      *resty.Response
}

//type HttpError struct {
//	HttpStatusCode int
//	Msg            string
//	Error          error
//	Errors         []*HttpError //重试逻辑的错误列表
//}

type res struct{}

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *Request) SetMethod(method string) {
	r.Method = method
}

func (r *Request) SetBody(body interface{}) {
	r.Body = body
}

func (r *Request) SetTimeout(timeout time.Duration) {
	r.Timeout = timeout
}

func (r *Request) SetRetryCount(tryCount int) {
	r.RetryCount = tryCount
}

func (r *Request) SetRetryWaitTime(waitTime time.Duration) {
	r.RetryWaitTime = waitTime
}

func (r *Request) SetHeader(name string, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[name] = value
}

func (r *Request) SetQueryParams(name string, value string) {
	if r.QueryParams == nil {
		r.QueryParams = make(map[string]string)
	}
	r.QueryParams[name] = value
}

func (r *Request) setUrl(url string) {
	r.Url = url
}

// 重复会去重
func (r *Request) AddUrl(url string) {
	if !goutils.IsInArray[string](r.Urls, url) {
		r.Urls = append(r.Urls, url)
	}
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
