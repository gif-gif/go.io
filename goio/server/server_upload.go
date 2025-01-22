package goserver

import (
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gin-gonic/gin"

	"io/ioutil"
	"os"
	"path"
)

type LocalUpload struct {
}

func (lu LocalUpload) Upload(c *gin.Context, uploadDir string) *Response {
	f, fh, err := c.Request.FormFile("file")
	if err != nil {
		return Error(7001, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return Error(7002, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}

	f.Close()

	md5str := goutils.MD5(data)
	filepath := md5str[0:2] + "/" + md5str[2:4] + "/"

	if err := os.MkdirAll(uploadDir+filepath, 0755); err != nil {
		return Error(7003, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}

	fileExt := path.Ext(fh.Filename)
	fileBasename := path.Base(fh.Filename)
	filename := filepath + fileBasename + "_" + md5str[8:16] + fileExt

	ff, err := os.Create(uploadDir + filename)
	if err != nil {
		return Error(7004, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}
	defer ff.Close()

	if _, err := ff.Write(data); err != nil {
		return Error(7005, fmt.Sprintf("上传失败，原因：%s", err.Error()))
	}

	return SuccessResponse(gin.H{
		"url": filename,
	})
}
