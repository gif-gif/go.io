package goserver

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/util/gconv"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
)

// 唯一ID
func RequestId(c *gin.Context) string {
	if v := c.GetHeader("X-Request-Id"); v != "" {
		return v
	}
	if v := c.Query("request_id"); v != "" {
		return v
	}
	if v := c.GetHeader("X-Trace-Id"); v != "" {
		return v
	}
	if v := c.Query("trace_id"); v != "" {
		return v
	}
	return uuid.New().String()
}

// 客户端IP
func ClientIP(c *gin.Context) string {
	if v := c.GetHeader("X-Real-IP"); v != "" {
		return v
	}
	if v := c.GetHeader("X-Forwarded-For"); v != "" {
		return v
	}
	if v := c.ClientIP(); v == "::1" {
		return "127.0.0.1"
	}
	return ""
}

// 请求数据
func RequestBody(c *gin.Context) interface{} {
	var (
		b           []byte
		buf         bytes.Buffer
		contentType = c.ContentType()
	)

	switch contentType {
	case "application/x-www-form-urlencoded", "text/xml", "application/json":
		io.Copy(&buf, c.Request.Body)
		b = buf.Bytes()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	default:
		return nil
	}

	if contentType == "application/json" {
		var body interface{}
		if err := json.Unmarshal(b, &body); err == nil {
			return body
		}
	}

	return string(b)
}

// 请求成功
func HttpSuccess(data interface{}) *Response {
	return SuccessResponse(data)
}

// 请求格式错误，比如参数格式、参数字段名等 不正确
func HttpBadRequest(msg string, showType ...uint32) *Response {
	return ErrorResponseX(gconv.String(400), msg, showType...)
}

// 用户没有访问权限，需要进行身份认证
func HttpUnauthorized(msg string, showType ...uint32) *Response {
	return ErrorResponseX(gconv.String(401), msg, showType...)
}

// 用户已进行身份认证，但权限不够
func HttpForbidden(msg string, showType ...uint32) *Response {
	return ErrorResponseX(gconv.String(403), msg, showType...)
}

// 接口不存在
func HttpNotFound(msg string, showType ...uint32) *Response {
	return ErrorResponseX(gconv.String(404), msg, showType...)
}

// 服务器内部错误
func HttpServerError(msg string, showType ...uint32) *Response {
	return ErrorResponseX(gconv.String(500), msg, showType...)
}

// 如需返回特殊错误码，调用此接口
func HttpFailForCode(code int64, msg string, showType ...uint32) *Response {
	return ErrorResponseX(gconv.String(code), msg, showType...)
}
