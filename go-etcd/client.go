package goetcd

import (
	"context"
	"crypto/tls"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
)

type GoEtcdClient struct {
	*clientv3.Client

	ctx  context.Context
	conf Config
}

func New(conf Config) (cli *GoEtcdClient, err error) {
	cli = &GoEtcdClient{ctx: context.TODO(), conf: conf}

	defer func() {
		if err != nil {
			time.Sleep(5 * time.Second)
			cli, _ = New(conf)
		}
	}()

	cfg := clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: 5 * time.Second,
		Logger:      zap.NewNop(),
	}

	if conf.Username != "" {
		cfg.Username = conf.Username
	}
	if conf.Password != "" {
		cfg.Password = conf.Password
	}

	if conf.TLS != nil {
		tlsInfo := &transport.TLSInfo{
			CertFile:      conf.TLS.CertFile,
			KeyFile:       conf.TLS.KeyFile,
			TrustedCAFile: conf.TLS.CAFile,
		}
		var clientConfig *tls.Config
		clientConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			golog.WithTag("go-etcd").WithField("config", conf).Error(err.Error())
			return
		}
		cfg.TLS = clientConfig
	}

	cli.Client, err = clientv3.New(cfg)
	if err != nil {
		golog.WithTag("go-etcd").WithField("config", conf).Error(err.Error())
	}

	return
}

// set key-value
func (cli *GoEtcdClient) Set(key, val string, opts ...clientv3.OpOption) (resp *clientv3.PutResponse, err error) {
	resp, err = cli.Client.Put(cli.ctx, key, val, opts...)
	if err != nil {
		golog.WithTag("go-etcd").WithField("key", key).WithField("val", val).Error(err)
	}
	return
}

// set key-value and return previous key-value
func (cli *GoEtcdClient) SetWithPrevKV(key, val string) (resp *clientv3.PutResponse, err error) {
	return cli.Set(key, val, clientv3.WithPrevKV())
}

// set key-value-ttl
func (cli *GoEtcdClient) SetTTL(key, val string, ttl int64, opts ...clientv3.OpOption) (resp *clientv3.PutResponse, err error) {
	if ttl == 0 {
		return cli.Set(key, val, opts...)
	}

	var lease *clientv3.LeaseGrantResponse

	lease, err = cli.Client.Grant(cli.ctx, ttl)
	if err != nil {
		golog.WithTag("go-etcd").WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err)
		return
	}

	_, err = cli.Client.Put(cli.ctx, key, val, clientv3.WithLease(lease.ID))
	if err != nil {
		golog.WithTag("go-etcd").WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err)
		return
	}

	return
}

// set key-value-ttl and return previous key-value
func (cli *GoEtcdClient) SetTTLWithPrevKV(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	return cli.SetTTL(key, val, ttl, clientv3.WithPrevKV())
}

// get value by key
func (cli *GoEtcdClient) Get(key string, opts ...clientv3.OpOption) (resp *clientv3.GetResponse, err error) {
	resp, err = cli.Client.Get(cli.ctx, key, opts...)
	if err != nil {
		golog.WithTag("go-etcd").WithField("key", key).Error(err)
	}
	return
}

// get string value by prefix key
func (cli *GoEtcdClient) GetString(key string) string {
	resp, err := cli.Get(key, clientv3.WithPrefix())
	if err != nil {
		return ""
	}
	if l := len(resp.Kvs); l == 0 {
		return ""
	}
	return string(resp.Kvs[0].Value)
}

// get array value by prefix key
func (cli *GoEtcdClient) GetArray(key string) (data []string) {
	data = []string{}

	resp, err := cli.Get(key, clientv3.WithPrefix())
	if err != nil {
		golog.WithTag("go-etcd").WithField("key", key).Error(err)
		return
	}

	for _, i := range resp.Kvs {
		data = append(data, string(i.Value))
	}

	return
}

// get map value by prefix key
func (cli *GoEtcdClient) GetMap(key string) (data map[string]string) {
	data = map[string]string{}

	resp, err := cli.Get(key, clientv3.WithPrefix())
	if err != nil {
		golog.WithTag("go-etcd").WithField("key", key).Error(err)
		return
	}

	for _, i := range resp.Kvs {
		key := string(i.Key)
		data[key] = string(i.Value)
	}

	return
}

// del key and return previous key-value
func (cli *GoEtcdClient) Del(key string, opts ...clientv3.OpOption) (resp *clientv3.DeleteResponse, err error) {
	opts = append(opts, clientv3.WithPrevKV())
	resp, err = cli.Client.Delete(cli.ctx, key, opts...)
	if err != nil {
		golog.WithTag("go-etcd").WithField("key", key).Error(err)
	}
	return
}

// del prefix key and return previous key-value
func (cli *GoEtcdClient) DelWithPrefix(key string) (resp *clientv3.DeleteResponse, err error) {
	return cli.Del(key, clientv3.WithPrefix())
}

// register service and keepalive
func (cli *GoEtcdClient) RegisterService(serviceName, addr string) (err error) {
	defer func() {
		if cli.Client == nil || err != nil {
			time.Sleep(3 * time.Second)
			cli.RegisterService(serviceName, addr)
		}
	}()

	var (
		ttl   int64 = 5
		em    endpoints.Manager
		lease *clientv3.LeaseGrantResponse
		ch    <-chan *clientv3.LeaseKeepAliveResponse
	)

	lease, err = cli.Client.Grant(cli.ctx, ttl)
	if err != nil {
		golog.WithTag("go-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	em, err = endpoints.NewManager(cli.Client, serviceName)
	if err != nil {
		golog.WithTag("go-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	serviceKey := serviceName + "/" + strconv.Itoa(int(lease.ID))
	err = em.AddEndpoint(cli.ctx, serviceKey, endpoints.Endpoint{Addr: addr}, clientv3.WithLease(lease.ID))
	if err != nil {
		golog.WithTag("go-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	ch, err = cli.Client.KeepAlive(cli.ctx, lease.ID)
	if err != nil {
		golog.WithTag("go-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	golog.WithTag("go-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Debug("服务注册成功")

	go func() {
		for {
			select {
			case <-cli.ctx.Done():
				golog.WithTag("go-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Warn("服务退出,收回注册信息")
				cli.Client.Revoke(cli.ctx, lease.ID)
				return

			case rsp := <-ch:
				if rsp == nil {
					golog.WithTag("go-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error("服务注册续租失效")
					cli.RegisterService(serviceName, addr)
					return
				}
			}
		}
	}()

	return
}

func (cli *GoEtcdClient) WatchWithContext(ctx context.Context, key string, withPrefix bool) <-chan []string {
	ch := make(chan []string, runtime.NumCPU()*2)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			err := cli.runWatchSession(ctx, key, withPrefix, ch)

			if err != nil {
				if ctx.Err() != nil {
					return
				}
				golog.WithTag("go-etcd").ErrorF("Etcd watch session failed: %v, retrying in 1s...", err)
				time.Sleep(1 * time.Second)
			} else {
				golog.WithTag("go-etcd").Info("Etcd watch session finished normally, reconnecting...")
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	return ch
}

func (cli *GoEtcdClient) runWatchSession(ctx context.Context, key string, withPrefix bool, ch chan<- []string) error {
	var baseOpts []clientv3.OpOption
	if withPrefix {
		baseOpts = append(baseOpts, clientv3.WithPrefix())
	}

	resp, err := cli.Client.Get(ctx, key, baseOpts...)
	if err != nil {
		return err
	}

	data := make(map[string]string)
	for _, kv := range resp.Kvs {
		data[string(kv.Key)] = string(kv.Value)
	}

	select {
	case ch <- cli.map2array(data):
	case <-ctx.Done():
		return ctx.Err()
	}

	watchStartRevision := resp.Header.Revision + 1
	watchOpts := append(baseOpts, clientv3.WithRev(watchStartRevision), clientv3.WithProgressNotify())
	wc := cli.Client.Watch(ctx, key, watchOpts...)

	for {
		select {
		case resp, ok := <-wc:
			if !ok {
				return nil
			}

			if resp.Err() != nil {
				return resp.Err()
			}

			if len(resp.Events) == 0 {
				continue
			}

			dirty := false
			for _, ev := range resp.Events {
				k := string(ev.Kv.Key)
				v := string(ev.Kv.Value)

				switch ev.Type {
				case clientv3.EventTypePut:
					if oldV, exists := data[k]; !exists || oldV != v {
						data[k] = v
						dirty = true
					}
				case clientv3.EventTypeDelete:
					if _, exists := data[k]; exists {
						delete(data, k)
						dirty = true
					}
				}
			}

			if dirty {
				select {
				case ch <- cli.map2array(data):
				case <-ctx.Done():
					return ctx.Err()
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// watch the key
func (cli *GoEtcdClient) Watch(key string) <-chan []string {
	return cli.WatchWithContext(cli.ctx, key, true)
}

func (cli *GoEtcdClient) map2array(data map[string]string) []string {
	var arrData []string
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	slices.SortStableFunc(keys, func(a, b string) int {
		return strings.Compare(a, b)
	})
	for _, k := range keys {
		arrData = append(arrData, data[k])
	}
	return arrData
}

// 检测 rpc服务是否启动
func (cli *GoEtcdClient) CheckRpcServices(rpcServices []string) bool {
	rpcStarted := true
	for _, service := range rpcServices {
		rpc := GetMap(service)
		if len(rpc) == 0 {
			if goio.Env == goio.TEST || goio.Env == goio.PRODUCTION {
				golog.WithTag("CheckRpcService").Fatal(service + " rpc is not started")
			} else {
				golog.WithTag("CheckRpcService").Error(service + " rpc is not started")
			}
			rpcStarted = false
		} else {
			golog.WithTag("CheckRpcService").Info(service + " rpc is started")
		}
	}

	return rpcStarted
}
