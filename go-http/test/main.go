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
	testRaceSpeed()
}

func testRequest() {
	req := gohttpx.Request{
		Url: "http://localhost:100",
		Urls: []string{
			"http://localhost:200",
			"http://localhost:300",
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
	err := gohttpx.HttpPostJson[gohttpx.Response](&req, res)
	if err != nil {
		//golog.ErrorF("Error: \n", err.ErrorInfo())
		fmt.Println(err.ErrorInfo())
	} else {
		fmt.Println(res)
	}

	time.Sleep(10 * time.Second)
}

func testRaceSpeed() {
	req := gohttpx.Request{
		IsConcurrency: true,
		IsAll:         true,
		Method:        gohttpx.POST,
		Url:           "http://localhost:100",
		Urls: []string{
			"http://localhost:20122/api/jump/account/check",
			"https://jumpjump.io/api/jump/account/check",
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
	err := gohttpx.HttpPostJson[gohttpx.Response](&req, res)
	if err != nil {
		golog.ErrorF("Error: \n", err.ErrorInfo())
	} else {
		golog.InfoF("res: \n", res)
	}

	time.Sleep(10 * time.Second)
}
