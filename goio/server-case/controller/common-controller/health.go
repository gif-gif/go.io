package common_controller

import (
	"github.com/gif-gif/go.io/goio/server"
	"github.com/gin-gonic/gin"
)

type Health struct {
}

func (this Health) DoHandle(ctx *gin.Context) *goserver.Response {
	ctx.String(200, "ok")
	return nil
}
