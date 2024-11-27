# go-redis

```
redis 连接和操作 基于 github.com/go-redis/redis 库
```
- 使用方法
```go
	config := goredis.Config{
		Name:     "goredis",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "goredis",
		AutoPing: true,
	}

	err := goredis.Init(config)
	if err != nil {
		golog.WithTag("goredis").Error(err)
	}

	cmd := goredis.Default().Set("goredis", "goredis")
	if cmd.Err() != nil {
		golog.WithTag("goredis").Error(cmd.Err())
	}
	v := goredis.Default().Get("goredis").Val()
	golog.WithTag("goredis").InfoF(v)
```
