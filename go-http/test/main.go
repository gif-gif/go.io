package main

import (
	"fmt"
	gohttpx "github.com/gif-gif/go.io/go-http/go-httpex"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)

	req := gohttpx.Request{
		Url: "http://localhost:100",
		Urls: []string{
			"http://localhost:200",
			"http://localhost:200",
			"http://localhost:400",
		},
		QueryParams: map[string]string{"name": "jk"},
		Timeout:     time.Second * 2,
	}
	type httpRequest struct {
		Email string `json:"email"`
	}

	req.Body = &httpRequest{
		Email: "test@gmail.com",
	}

	res := &gohttpx.Response{}
	res, err := gohttpx.HttpPostJson[gohttpx.Response](req, res)
	if err != nil {
		golog.ErrorF("Error: %+v\n", err)
	} else {
		fmt.Println(res)
	}

	time.Sleep(10 * time.Second)
}
