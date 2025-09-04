# GoPool for workers 并发请求 go routines 池

- 基于 https://github.com/panjf2000/ants 封装

Examples
```
	pool, _ := gopool.New(100, ants.WithPreAlloc(true))
	err := pool.Submit(func() {

	})
	if err != nil {
		golog.ErrorF("Submit failed: %v", err)
	}
```
