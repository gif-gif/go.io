package goserver

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 定义控制器抽象类
type iController interface {
	DoHandle(c *gin.Context) *Response
}

// 定义控制器调用实现
func Handler(controller iController) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := controller.DoHandle(c)

		if resp == nil {
			return
		}

		c.Set("__response", resp.Copy())

		// 计算执行时间
		beginTime := c.GetTime("__begin_time")
		if !beginTime.IsZero() {
			c.Header("X-Response-Duration", fmt.Sprintf("%dms", time.Since(beginTime)/1e6))
		}

		if !defaultOptions.encryptionEnable {
			c.JSON(200, resp)
			return
		}

		switch strings.ToUpper(c.Request.Method) {
		case "POST", "PUT":
		default:
			c.JSON(200, resp)
			return
		}

		switch strings.ToLower(c.Request.Header.Get("Content-Type")) {
		case "multipart/form-data":
			c.JSON(200, resp)
			return
		}

		if _, ok := defaultOptions.encryptionExcludeUris[c.Request.RequestURI]; ok {
			c.JSON(200, resp)
			return
		}

		b, err := json.Marshal(&resp.Data)
		if err != nil {
			c.JSON(500, Error(5003, "数据解析失败，原因："+err.Error()))
			return
		}

		body, err := defaultOptions.encryption.Encode(b)
		if err != nil {
			c.JSON(500, Error(5004, "数据解析失败，原因："+err.Error()))
			return
		}

		resp.Data = body
		c.JSON(200, resp)

		return
	}
}
