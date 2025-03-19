package goi

import (
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
)

// 单机
func Redis(names ...string) *goredis.GoRedis {
	return goredis.GetClient(names...)
}

// 集群
func Redisc(names ...string) *goredisc.GoRedisC {
	return goredisc.GetClient(names...)
}
