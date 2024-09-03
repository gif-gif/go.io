# gohttpx 模块
- 快捷使用http 请求
```go
resp, err := gohttpx.PostJson(Email.Api, data)
if err != nil {
    logx.Error("邮件发送失败：" + string(data) + ",error:" + err.Error())
    return err
} else {
    logx.Infof("executeSendEmail result:%s data:%s", string(resp), string(data))
}
return nil
```

```go
rst, err := gohttpx.Post(url, []byte(postParams))
```
```go
clientIp := gohttpx.GetClientIp(request)
```
