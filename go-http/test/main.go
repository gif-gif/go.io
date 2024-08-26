package main

import (
	"context"
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-http"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	logx.DisableStat()
	testRaceSpeed()
	<-gocontext.Cancel().Done()
}

func testRequest() {
	type httpRequest struct {
		Email string `json:"email"`
	}

	req := &gohttp.Request{
		Url: "/main",
		Urls: []string{
			"/main1",
			"/main2",
			"/main3",
		},
		QueryParams: map[string]string{"name": "jk"},
		Timeout:     time.Second * 2,
		Body: &httpRequest{
			Email: "test@gmail.com",
		},
	}

	gh := &gohttp.GoHttp[gohttp.Response]{
		Request: req,
		BaseUrl: "http://localhost",
		Headers: map[string]string{
			"User-Agent": "github.com/gif-gif/go.io",
		},
	}

	rst, err := gh.HttpPostJson(context.Background())
	if err != nil {
		golog.WithTag("http").Error(err.Error())
	} else {
		fmt.Println(rst)
	}
}

func testRaceSpeed() {
	type httpRequest struct {
		Email string `json:"email"`
	}
	req := gohttp.Request{
		Method: gohttp.POST,
		Urls: []string{
			"https://jumpjump.io/api/jump/account/check",
			"http://localhost:400",
		},
		QueryParams: map[string]string{"name": "jk"},
		Timeout:     time.Second * 2,
	}

	req.Body = &httpRequest{
		Email: "test111@gmail.com",
	}

	gh := &gohttp.GoHttp[gohttp.Response]{
		Request: &req,
	}

	rst, err := gh.HttpConcurrencyRequest()
	if err != nil {
		golog.ErrorF("Error: \n", err.Error())
	} else {
		golog.InfoF("res: \n", rst.Code, rst.Msg)
	}
}
