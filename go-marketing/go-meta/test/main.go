package main

import (
	"context"
	gocontext "github.com/gif-gif/go.io/go-context"
	golog "github.com/gif-gif/go.io/go-log"
	gometa "github.com/gif-gif/go.io/go-marketing/go-meta"
	"github.com/gif-gif/go.io/goio"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	err := gometa.Init(context.Background(), gometa.Config{
		AccessToken: "token",
		StartDate:   "2024-01-01",
		EndDate:     "2024-01-01",
		PageSize:    200,
	})

	if err != nil {
		golog.WithTag("goMeta").Error(err.Error())
		return
	}
	//gometa.Default().AccessKeys("123")
	//if err != nil {
	//	golog.WithTag("goMeta").Error(err.Error())
	//	return
	//}

	//golog.WithTag("goMeta").Info(res)

	<-gocontext.Cancel().Done()
}
