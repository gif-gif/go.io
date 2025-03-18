package goredisc

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	Name string `yaml:"Name" json:"name,optional"`
	//NamePrefix string `yaml:"NamePrefix" json:"namePrefix,optional"`
	Addrs    []string `yaml:"Addrs" json:"addrs,optional"`
	Password string   `yaml:"Password" json:"password,optional"`
	DB       int      `yaml:"DB" json:"db,optional"`
	Prefix   string   `yaml:"Prefix" json:"prefix,optional"`
	AutoPing bool     `yaml:"AutoPing" json:"autoPing,optional"`
	TLS      bool     `yaml:"TLS" json:"tls,optional"`
	Type     string   `yaml:"Type" json:",default=node,options=node|cluster"`
	Weight   int      `yaml:"Weight" json:",default=100"` //for gozero TODO: 这里需要有多个节点的配置
}

type ClusterConf Config

// 兼容go-zore
func (c ClusterConf) GetCacheConf() cache.CacheConf {
	cacheConf := make([]cache.NodeConf, 0, len(c.Addrs))
	for _, addr := range c.Addrs {
		cacheConf = append(cacheConf, cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: addr,
				Pass: c.Password,
				Tls:  c.TLS,
				Type: c.Type,
			},
			Weight: c.Weight,
		})
	}
	return cacheConf
}
