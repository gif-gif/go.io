# Meta(Facebook) Marketing Api
```go
meta := gometa.Market{
    BaseApi:     "https://graph.facebook.com/v17.0",
    AccessToken: "token",
    StartDate:   "2024-01-01",
    EndDate:     "2024-01-01",
    PageSize:    200,
}
res, err := meta.GetAccountsByBusinessId("15738715864408601")
if err != nil {
    golog.WithTag("goMeta").Error(err.Error())
}

golog.WithTag("goMeta").Info(res)
```