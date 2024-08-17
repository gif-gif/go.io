package main

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
)

func main() {
	// 加密
	key, iv, err := goutils.GenerateAESKeyAndIV()
	if err != nil {
		golog.WithTag("aes").Error(err)
		return
	}

	fmt.Printf("密钥: %x\n", key)
	fmt.Printf("IV: %x\n", iv)

	fmt.Println(key, iv)
	golog.WithTag("aes1111111").Info(key, iv)
	key, err = goutils.GenerateAESKey()
	if err != nil {
		golog.WithTag("aes").Error(err)
		return
	}
	golog.WithTag("aes22222222").Info(key)
}
