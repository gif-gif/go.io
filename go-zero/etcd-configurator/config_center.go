package etcd_configurator

import (
	goetcd "github.com/gif-gif/go.io/go-etcd"
	golog "github.com/gif-gif/go.io/go-log"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/configcenter/subscriber"
)

func NewConfigCenter[T any](key string, etcd goetcd.Config, callback func(t T)) (configurator.Configurator[T], error) {
	ss := subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
		Hosts: etcd.Endpoints, // etcd 地址
		User:  etcd.Username,
		Pass:  etcd.Password,
		Key:   key, // 配置key
	})

	// 创建 configurator
	cc := configurator.MustNewConfigCenter[T](configurator.Config{
		Type: "json", // 配置值类型：json,yaml,toml
	}, ss)

	// 获取配置
	// 注意: 配置如果发生变更，调用的结果永远获取到最新的配置
	v, err := cc.GetConfig()
	if err != nil {
		return nil, err
	}
	callback(v)
	// 如果想监听配置变化，可以添加 listener
	cc.AddListener(func() {
		v, err := cc.GetConfig()
		if err != nil {
			golog.WithTag("etcd").Fatal(err)
			return
		}
		callback(v)
	})

	return cc, nil
}
