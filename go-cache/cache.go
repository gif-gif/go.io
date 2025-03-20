package gocache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"time"
)

type GoCache struct {
	SharedCache *cache.Cache
	ctx         context.Context
	Config      Config
}

func New(conf Config) (cli *GoCache) {
	cli = &GoCache{ctx: context.TODO(), Config: conf}
	cli.SharedCache = cache.New(time.Second*time.Duration(conf.DefaultExpiration), time.Duration(conf.CleanupInterval)*time.Second)
	return
}

func (cli *GoCache) Get(key string) (interface{}, bool) {
	if cli.SharedCache == nil {
		return nil, false
	}
	return cli.SharedCache.Get(key)
}

func (cli *GoCache) Set(key string, data interface{}, duration time.Duration) {
	if cli.SharedCache == nil {
		return
	}
	cli.SharedCache.Set(key, data, duration)
}
