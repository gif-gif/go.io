package common_controller

import (
	gofile "github.com/gif-gif/go.io/go-file"
	goserver "github.com/gif-gif/go.io/goio/server"
	"github.com/gin-gonic/gin"
)

type FileDownload struct {
}

//
//func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
//	ds := gofile.NewGoDownload(ctx, "test.csv", ctx.Writer, ctx.Request)
//	go ds.Run()
//	file := "/Users/Jerry/Desktop/test/ios_file.csv"
//	err := ds.OutputByLine(file)
//	if err != nil {
//		return nil
//	}
//	ds.WaitDone()
//	return nil
//}

func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
	ds := gofile.NewGoDownload(ctx, "test.csv", ctx.Writer, ctx.Request)
	go ds.Run()
	filePath := "/Users/Jerry/Desktop/test/ios_file3.csv"
	err := gofile.ReadLines(filePath, func(chunk string) error {
		ds.Write([]byte(chunk + "\n"))
		return nil
	})

	ds.Close()
	if err != nil {
		ds.Error(err)
		return nil
	}
	ds.WaitDone()
	return nil
}
