# go.io
framework interfaces and configurations for libraries

```go
package main

import (
	"github.com/gif-gif/go.io/goio"
	"github.com/gif-gif/go.io/goio/test/router"
)

func main() {
	startSever()
}

func startSever() {
	goio.Init(goio.DEVELOPMENT)
	s := goio.NewServer(
		goio.ServerNameOption("serverName"),
		goio.EnvOption(goio.Env),
		//goio.EnableEncryptionOption(),
		goio.PProfEnableOption(false),
		goio.NoLogPathsOption("/captcha/get"),
	)

	// 路由
	router.Routes(s.Group("/"))
	// 启动
	s.Run(":1000")
}
```

```javascript
curl --location --request POST 'http://127.0.0.1:1000/health' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)'
```
