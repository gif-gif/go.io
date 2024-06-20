package main

import (
	golog "github.com/gif-gif/go.io/go-log"
	gometa "github.com/gif-gif/go.io/go-marketing/go-meta"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	meta := gometa.Market{
		BaseApi: "https://graph.facebook.com/v17.0",
	}
	res, err := meta.GetAccountsByBusinessId("15738715864408601")
	if err != nil {
		golog.WithTag("goMeta").Error(err.Error())
	}

	golog.WithTag("goMeta").Info(res)
	time.Sleep(10000 * time.Second)
}
