package common_controller

import (
	"github.com/gif-gif/go.io/goio"
	"github.com/gin-gonic/gin"
)

type Health struct {
}

func (this Health) DoHandle(ctx *gin.Context) *goio.Response {
	ctx.String(200, "ok")
	return nil
}
