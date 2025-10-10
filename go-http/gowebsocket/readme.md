# GoWebsocket 

## Server 
```go
// server.go
package main

import (
	"log"

	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-http/gowebsocket"
)

func main() {
	gowebsocket.InitServer(gowebsocket.ServerConfig{
		Port: 8080,
	})

	if err := gowebsocket.DefaultServer().Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

	<-gocontext.Cancel().Done()
}

```

## Client
```go
// client.go
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"

	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-http/gowebsocket"
)

func main() {
	var addr = flag.String("addr", "localhost", "服务器地址")
	var clientID = flag.String("id", "", "客户端ID")
	flag.Parse()

	// 如果没有指定客户端ID，则自动生成
	if *clientID == "" {
		*clientID = fmt.Sprintf("client_%d", rand.Int63())
	}

	gowebsocket.InitClient(gowebsocket.ClientConfig{
		ClientID:          *clientID,
		Addr:              *addr,
		Port:              8080,
		HeartBeatInterval: 10,
	})

	// 连接到服务器
	if err := gowebsocket.DefaultClient().Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	log.Printf("[%s] 客户端启动", *clientID)
	log.Printf("[%s] 心跳间隔: %v", *clientID, gowebsocket.DefaultClient().ClientConfig.HeartBeatInterval)

	// 启动客户端
	gowebsocket.DefaultClient().Start()

	log.Printf("[%s] 客户端已退出", *clientID)
	<-gocontext.Cancel().Done()
}

```