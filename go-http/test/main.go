package main

import (
	"context"
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-http"
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
	gh := &gohttp.GoHttp[gohttp.Response]{}
	gh.SetBaseUrl("http://localhost")
	gh.AddGlobalHeader("User-Id", "123")
	req := gohttp.Request{
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

	rst, err := gh.HttpPostJson(context.Background(), &req)
	if err != nil {
		golog.WithTag("http").Error(err.Error())
	} else {
		fmt.Println(rst)
	}

	time.Sleep(10 * time.Second)
}

func testRaceSpeed() {
	gh := &gohttp.GoHttp[gohttp.Response]{}

	req := gohttp.Request{
		Method: gohttp.POST,
		Urls: []string{
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
		Email: "test111@gmail.com",
	}

	rst, err := gh.HttpConcurrencyRequest(&req)
	if err != nil {
		golog.ErrorF("Error: \n", err.Error())
	} else {
		golog.InfoF("res: \n", rst.Code, rst.Msg)
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
