package main

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	gopool "github.com/gif-gif/go.io/go-pool"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio"
	"github.com/gogf/gf/util/gconv"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	testSign()
	gp := gopool.NewFixedSizePool(10, 10)
	defer gp.StopAndWait()
	s := []string{"s1", "s2", "s3"}
	sa := goutils.NewSafeSlice[string]()
	sa.Sets(s)
	for i := 0; i < 10; i++ { //并发
		gp.Submit(func() {
			fmt.Println(sa.Get())
		})
	}
	time.Sleep(30 * time.Second)
}

func testSign() {
	ts := time.Now().Unix()
	sign := goutils.Md5([]byte(gconv.String(ts) + "123456"))
	golog.WithTag("sign").Info(ts, sign)
}
