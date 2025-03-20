package main

import (
	"fmt"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
)

func main() {
	c := goredisc.ClusterConf{
		Name: "goredis",
	}

	config := c.GetConfig()

	fmt.Println(config)
}
