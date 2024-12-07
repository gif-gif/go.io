package common_controller

import (
	gofile "github.com/gif-gif/go.io/go-file"
	goserver "github.com/gif-gif/go.io/goio/server"
	"github.com/gin-gonic/gin"
)

type FileDownload struct {
}

//func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
//	ds := gofile.NewGoDownload(ctx, "test.csv", ctx.Writer, ctx.Request)
//	err := ds.SetFileHeaders()
//	file := "/Users/Jerry/Desktop/test/ios_file3.csv"
//	err = ds.Output(file)
//	if err != nil {
//		http.Error(ctx.Writer, "Streaming unsupported!", http.StatusInternalServerError)
//		return nil
//	}
//	return nil
//}

func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
	ds := gofile.NewGoDownload(ctx, "test.csv", ctx.Writer, ctx.Request)
	err := ds.SetFileHeaders()
	if err != nil {
		return nil
	}

	filePath := "/Users/Jerry/Desktop/test/ios_file3.csv"
	err = gofile.ReadLines(filePath, func(chunk string) error {
		err = ds.Write([]byte(chunk + "\n"))
		return err
	})

	if err != nil {
		ds.Error(err)
		return nil
	}
	return nil
}
