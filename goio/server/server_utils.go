package goserver

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
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
