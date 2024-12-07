# 文件操作相关模块

- 文件下载功能

### 方法 1
```go
func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
    ds := gofile.NewGoDownload(ctx, "downloaded_file.csv", ctx.Writer, ctx.Request)
    go ds.Run()
    filePath := "file.csv"
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
```

### 方法 2
```go
func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
	ds := gofile.NewGoDownload(ctx, "downloaded_file.csv", ctx.Writer, ctx.Request)
	go ds.Run()
	file := "file.csv"
	err := ds.OutputByLine(file)
	if err != nil {
		return nil
	}
	ds.WaitDone()
	return nil
}
```

