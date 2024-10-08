package common_controller

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio/server"
	"github.com/gin-gonic/gin"
)

type Captcha struct {
	request struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
}

func (this Captcha) DoHandle(ctx *gin.Context) *goserver.Response {
	if err := ctx.ShouldBind(&this.request); err != nil {
		return goserver.Error(7001, err.Error())
	}
	rsp := goutils.CaptchaGet(this.request.Width, this.request.Height)
	return goserver.Success(rsp)
}
