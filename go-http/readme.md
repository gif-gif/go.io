# Http 请求封装

```go
package main

import (
	"fmt"
	gohttp "github.com/gif-gif/go.io/go-http"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	gh := &gohttp.GoHttp[gohttp.Response]{}
	req := gohttp.Request{
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

	res,err := gh.HttpPostJson(&req)
	if err != nil {
		golog.ErrorF("Error: %+v\n", err)
	} else {
		fmt.Println(res)
	}

	time.Sleep(10 * time.Second)
}

```

- Custom Root Certificates and Client Certificates
- Custom Root Certificates and Client Certificates from string
- Save HTTP Response into File
