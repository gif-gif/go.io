package goi

import goredis "github.com/gif-gif/go.io/go-db/go-redis"

func Redis(names ...string) *goredis.GoRedis {
	return goredis.GetClient(names...)
}
