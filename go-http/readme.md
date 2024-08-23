# Http 请求封装

```go
package main

import (
	"context"
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-http"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
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

```

- Custom Root Certificates and Client Certificates
- Custom Root Certificates and Client Certificates from string
- Save HTTP Response into File
