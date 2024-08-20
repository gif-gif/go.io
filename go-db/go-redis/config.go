package goredis

import (
	redis0 "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	Name     string `yaml:"Name" json:"name,optional"`
	Addr     string `yaml:"Addr" json:"addr,optional"`
	Password string `yaml:"Password" json:"password,optional"`
	DB       int    `yaml:"DB" json:"db,optional"`
	Prefix   string `yaml:"Prefix" json:"prefix,optional"`
	AutoPing bool   `yaml:"AutoPing" json:"autoPing,optional"`
	TLS      bool   `yaml:"TLS" json:"tls,optional"`
}

type ClusterConf []Config

// 兼容go-zore
func (c ClusterConf) GetCacheConf() cache.CacheConf {
	cacheConf := make([]cache.NodeConf, 0, len(c))
	for _, conf := range c {
		cacheConf = append(cacheConf, cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: conf.Addr,
				Pass: conf.Password,
				Tls:  conf.TLS,
			},
		})
	}
	return cacheConf
}

type RedisHook struct {
	Prefix string `yaml:"Prefix" json:"prefix,optional"`
}

func (s *RedisHook) DialHook(next redis0.DialHook) redis0.DialHook {
	return next
}

func (s *RedisHook) ProcessHook(next redis0.ProcessHook) redis0.ProcessHook {
	return next
}

func (s *RedisHook) ProcessPipelineHook(next redis0.ProcessPipelineHook) redis0.ProcessPipelineHook {
	return next
}

//func (s *XzRedis) WrapKey(key string) string {
//	return s.KeyPrefix + ":" + key
//}
