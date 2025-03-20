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
