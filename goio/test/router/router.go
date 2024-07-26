package router

import (
	"github.com/gif-gif/go.io/goio"
	common_controller "github.com/gif-gif/go.io/goio/test/controller/common-controller"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.POST("/health", goio.Handler(common_controller.Health{}))

	r.Use(verifySign)
	{
		r.POST("/captcha/get", goio.Handler(common_controller.Captcha{}))
	}

}
