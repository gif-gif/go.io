# go-redis

```
redis 连接和操作 基于 github.com/go-redis/redis 库
```
- 建议优先使用 goredisc （兼容集群和单点）
- 使用方法
```go
	config := goredisc.Config{
		Name:     "goredis",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "goredis",
		AutoPing: true,
	}

	err := goredisc.Init(config)
	if err != nil {
		golog.WithTag("goredis").Error(err)
	}

	cmd := goredisc.Default().Set("goredis", "goredis")
	if cmd.Err() != nil {
		golog.WithTag("goredis").Error(cmd.Err())
	}
	v := goredisc.Default().Get("goredis").Val()
	golog.WithTag("goredis").InfoF(v)
```
