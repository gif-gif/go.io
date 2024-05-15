package goio

type res struct{}

type OrderItem struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc"`
}

type Page struct {
	PageNo    int64        `json:"page_no,optional"`
	PageSize  int64        `json:"page_size,optional"`
	StartTime int64        `json:"start_time,optional"`
	EndTime   int64        `json:"end_time,optional"`
	SortBy    []*OrderItem `json:"sort_by,optional"`
}

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
