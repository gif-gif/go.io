package main

import (
	"context"
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	gohttpx "github.com/gif-gif/go.io/go-http/go-httpex"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	testRaceSpeed()

	<-gocontext.Cancel().Done()
}

func testRequest() {
	gohttpx.SetBaseUrl("http://localhost")
	gohttpx.AddGlobalHeader("User-Id", "123")
	req := gohttpx.Request{
		Url: "/main",
		Urls: []string{
			"/main1",
			"/main2",
			"/main3",
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
	err := gohttpx.HttpPostJson[gohttpx.Response](context.Background(), &req, res)
	if err != nil {
		golog.WithTag("http").Error(err.ErrorInfo())
	} else {
		fmt.Println(res)
	}

	time.Sleep(10 * time.Second)
}

func testRaceSpeed() {
	req := gohttpx.Request{
		Method: gohttpx.POST,
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
	err := gohttpx.HttpConcurrencyRequest[gohttpx.Response](&req, res)
	if err != nil {
		golog.ErrorF("Error: \n", err.ErrorInfo())
	} else {
		golog.InfoF("res: \n", res)
	}

	time.Sleep(10 * time.Second)
}

func testChan() {
	c := make(chan int)

	goutils.AsyncFunc(func() {
		for {
			select {
			case _, ok := <-c:
				if !ok {
					fmt.Println("test")
				}
			}
		}

		for v := range c {
			fmt.Println("value:", v)
		}
	})

	close(c)

	fmt.Println("done")
}
