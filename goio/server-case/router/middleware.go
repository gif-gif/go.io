package router

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio"
	"github.com/gin-gonic/gin"
	"strings"
)

func verifySign(c *gin.Context) {
	ua := c.GetHeader("User-Agent")
	tss := c.GetHeader("X-Request-Timestamp")
	sign := c.GetHeader("X-Request-Sign")

	if strings.ToLower(sign) != strings.ToLower(goutils.SHA1([]byte(tss+goutils.MD5([]byte(ua))))) {
		c.AbortWithStatusJSON(403, goio.Error(40301, "签名错误"))
		return
	}

	c.Next()
}

func verifyCaptcha(c *gin.Context) {
	if 1 == 1 {
		c.AbortWithStatusJSON(403, goio.Error(40301, "验证码错误"))
		return
	}
	c.Next()
}
