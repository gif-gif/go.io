# go.io
- Golang Development Framework Continuously Developing and Updating

## How to install
```
go get -u github.com/gif-gif/go.io
```

## How to use

### Worker pool with dynamic size
```go
package main

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	gomessage "github.com/gif-gif/go.io/go-message"
	"github.com/gif-gif/go.io/goio"
)

type Config struct {
	Name   string
	FeiShu string
	Mode   string
}

func main() {
	c := Config{}
	goio.Init(goio.Environment(c.Mode))
	golog.WithHook(func(msg *golog.Message) {
		if msg.Level > golog.ERROR { //致命错误以上
			gomessage.FeiShu(c.FeiShu, fmt.Sprintf(">> %s/%s >> %s",
				c.Name, c.Mode, string(msg.JSON())))
		}
	})
}

```
- 各个模块功能在go-[模块]目录中readme或者testCase中找到使用方法

## 设计目标
- goio 提供了常用库封装，支持必要的简洁使用功能，在其之上可以进二次开发，以提供更好的代码维护；
- 以跨平台跨项目为首要原则，以减少二次开发的成本；
- 各个模块逻辑保持唯一不重复

## 开发规范
- dev 分之开发，跑测试case，确定没问题 合并到 main 分支跑测试case
- main 发布 release，版本号修改

### 对代码的修改
#### 功能性问题
- 请提交至少一个测试用例（Test Case）来验证对现有功能的改动。

#### 性能相关
- 请提交必要的测试数据来证明现有代码的性能缺陷，或是新增代码的性能提升。

#### 新功能
- 如果新增功能对已有功能不影响，请提供可以开启/关闭的开关（如 flag），并使新功能保持默认关闭的状态；
- 大型新功能（比如增加一个新的模块）开发之前，请先提交一个 issue，讨论完毕之后再进行开发。

## Thanks
- https://github.com/IBM/sarama
- https://gorm.io/gorm
- https://github.com/redis/go-redis
- https://github.com/aliyun/aliyun-oss-go-sdk
- https://github.com/go-co-op/gocron
- https://github.com/minio/minio-go
- https://github.com/oschwald/geoip2-golang
- https://github.com/ip2location/ip2location-go
- https://github.com/wechatpay-apiv3/wechatpay-go
- https://github.com/smartwalle/alipay
- https://github.com/360EntSecGroup-Skylar/excelize
- https://github.com/gin-gonic/gin
- https://github.com/go-resty/resty
- https://github.com/olivere/elastic
- https://github.com/mongodb/mongo-go-driver
- https://github.com/alitto/pond
