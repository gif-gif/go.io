package router

import (
	"github.com/gif-gif/go.io/goio/server"
	common_controller2 "github.com/gif-gif/go.io/goio/server-case/controller/common-controller"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.POST("/health", goserver.Handler(common_controller2.Health{}))
	r.POST("/login", goserver.Handler(common_controller2.Login{}))
	r.POST("/captcha/get", goserver.Handler(common_controller2.Captcha{}))
	r.POST("/file/download", goserver.Handler(common_controller2.FileDownload{}))

	r.Use(verifySign)
	{

	}

}
