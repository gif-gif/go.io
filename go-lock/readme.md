# Go-LOCK
- 用法
```go
package main

import (
	"fmt"
	golock "github.com/gif-gif/go.io/go-lock"
	gopool "github.com/gif-gif/go.io/go-pool"
	"github.com/gif-gif/go.io/goio"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	testGoPool()
	//testSyncLock()
}

func testGoPool() {
	var (
		count int
	)
	lock := golock.New()
	//并发池
	gp := gopool.NewFixedSizePool(10, 10)
	defer gp.StopAndWait()
	group := gp.NewTaskGroup()
	for i := 0; i < 10; i++ {
		group.Submit(func() {
			for i := 1000; i > 0; i-- {
				lock.WLockFunc(func(parameters ...any) {
					count++
				})
			}
			fmt.Println(count)
		})
	}
	group.Wait()
}

func testSyncLock() {
	var (
		count int
	)

	lock := golock.New()
	for i := 0; i < 2; i++ {
		go func() {
			for i := 1000; i > 0; i-- {
				//lock.WLock()
				//count++
				//lock.WUnlock()

				lock.WLockFunc(func(parameters ...any) {
					count++
				})
			}
			fmt.Println(count)
		}()
	}

	fmt.Scanf("\n") //等待子线程全部结束
}

```
### 分布式锁
- RedisLock 用 https://github.com/go-redsync/redsync
- EtcdLock 官方
```go
import "go.etcd.io/etcd/client/v3/concurrency"

// 官方推荐用法
session, _ := concurrency.NewSession(cli, concurrency.WithTTL(5))
defer session.Close()

mutex := concurrency.NewMutex(session, "/my-lock/")
if err := mutex.Lock(context.TODO()); err == nil {
    // 业务代码...
    mutex.Unlock(context.TODO())
}
```
