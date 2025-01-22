package common_controller

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio/server"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Request struct {
		Account     string `json:"account" binding:"required"`
		Password    string `json:"password" binding:"required"`
		CaptchaId   string `json:"captcha_id" binding:"required"`
		CaptchaCode string `json:"captcha_code" binding:"required"`
	}
	Response LoginResponse
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	IsSuper  bool   `json:"is_super"`
}

func (this Login) DoHandle(ctx *gin.Context) *goserver.Response {
	if err := ctx.ShouldBind(&this.Request); err != nil {
		return goserver.ErrorWithValidate(err, map[string]string{
			"account_required":      "登录账号 为空",
			"password_required":     "登录密码 为空",
			"captcha_id_required":   "验证码 无效",
			"captcha_code_required": "验证码 为空",
		})
	}

	if !goutils.CaptchaVerify(this.Request.CaptchaId, this.Request.CaptchaCode) {
		return goserver.Error(7001, "验证码错误")
	}

	var userId int64 = 1
	var userName = "tes"
	token, err := goserver.CreateToken("AppKey", userId)
	if err != nil {
		return goserver.Error(7004, "登录失败，提示信息："+err.Error())
	}

	this.Response = LoginResponse{
		Token:    token,
		Username: userName,
		IsSuper:  true,
	}

	if this.Response.Username == "" {
		this.Response.Username = userName
	}

	goutils.AsyncFunc(func() {

	})

	ctx.Set("user_id", userId)
	ctx.Set("username", userName)

	return goserver.SuccessResponse(this.Response)
}
