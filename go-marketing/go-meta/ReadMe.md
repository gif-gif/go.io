# Meta(Facebook) Marketing Api
```go
meta := gometa.Market{
    BaseApi: "https://graph.facebook.com/v17.0",
}
res, err := meta.GetAccountsByBusinessId("15738715864408601")
if err != nil {
    golog.WithTag("goMeta").Error(err.Error())
}
```