package router

import (
	"github.com/gif-gif/go.io/goio"
	common_controller2 "github.com/gif-gif/go.io/goio/server-case/controller/common-controller"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.POST("/health", goio.Handler(common_controller2.Health{}))
	r.POST("/login", goio.Handler(common_controller2.Login{}))
	r.POST("/captcha/get", goio.Handler(common_controller2.Captcha{}))
	r.Use(verifySign)
	{

	}

}
