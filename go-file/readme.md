# 文件操作相关模块

- 文件下载功能

### 方法 1
```go
func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
    ds := gofile.NewGoDownload(ctx, "downloaded_file.csv", ctx.Writer, ctx.Request)
    err := ds.SetFileHeaders()
    if err != nil {
        return nil
    }
    
    filePath := "file.csv"
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
```

### 方法 2
```go
func (this FileDownload) DoHandle(ctx *gin.Context) *goserver.Response {
	ds := gofile.NewGoDownload(ctx, "downloaded_file.csv", ctx.Writer, ctx.Request)
	err := ds.SetFileHeaders()
	file := "file.csv"
	err = ds.Output(file)
	if err != nil {
		http.Error(ctx.Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return nil
	}
	return nil
}
```

