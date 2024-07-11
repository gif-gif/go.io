package main

import (
	"fmt"
	golock "github.com/gif-gif/go.io/go-lock"
	gopool "github.com/gif-gif/go.io/go-pool"
	"github.com/gif-gif/go.io/goio"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	lock := golock.New()
	//并发池
	gp := gopool.NewFixedSizePool(100, 100)
	defer gp.StopAndWait()

	for i := 0; i < 10; i++ {
		gp.Submit(func() {
			lock.Lock(func(parameters ...any) {
				fmt.Println(parameters[0])
			}, i)
		})
	}
}
