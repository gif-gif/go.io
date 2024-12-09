
# go.io 
Golang Development Framework Continuously Developing and Updating

[![Stars](https://img.shields.io/github/stars/gif-gif/go.io)](https://github.com/gif-gif/go.io/stargazers)
[![Watchers](https://img.shields.io/github/watchers/gif-gif/go.io.svg?label=Watchers&style=social)](https://github.com/gif-gif/go.io/watchers)
[![Forks](https://img.shields.io/github/forks/gif-gif/go.io)](https://github.com/gif-gif/go.io/forks)
[![Golang](https://img.shields.io/badge/Go-%3E%3D%201.23-yellowgreen?style=flat)](https://github.com/gif-gif/go.io)
[![Release](https://img.shields.io/github/v/release/gif-gif/go.io.svg)](https://github.com/gif-gif/go.io/releases)
[![License](https://img.shields.io/github/license/gif-gif/go.io)](https://github.com/gif-gif/go.io?tab=Apache-2.0-1-ov-file)
[![Issues](https://img.shields.io/github/issues/gif-gif/go.io)](https://github.com/gif-gif/go.io/issues)
[![Commits](https://img.shields.io/github/commit-activity/m/gif-gif/go.io.svg?style=flat&label=commits)](https://github.com/gif-gif/go.io/graphs/commit-activity)
[![Downloads](https://img.shields.io/github/downloads/gif-gif/go.io/total.svg)](https://github.com/gif-gif/go.io/releases)
[![Contributors](https://img.shields.io/github/contributors/gif-gif/go.io)](https://github.com/gif-gif/go.io/graphs/contributors)
[![PullRequest](https://img.shields.io/github/issues-pr/gif-gif/go.io?color=117abe)](https://github.com/gif-gif/go.io/pulls)


## 设计目标
- goio 提供了常用库封装，支持必要的简洁使用功能，在其之上可以进二次开发，以提供更好的代码维护；
- 以跨平台跨项目为首要原则，以减少二次开发的成本；
- 各个模块逻辑保持唯一不重复，但模块之前相互便捷使用实现复杂逻辑开发

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


## How to install
```
go get -u github.com/gif-gif/go.io
```

## How to use

### 项目启动时初始化 ,日志文件默认在项目目录下 logs/date.log  
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
	//golog.SetAdapter(golog.NewFileAdapter()) //当前工程目录logs/date.log, 可通过这个设置改变日志输出目录 
	//不设置 golog.SetAdapter 默认控制台输出
	golog.WithHook(func(msg *golog.Message) {
		if msg.Level > golog.ERROR { //致命错误以上
			gomessage.FeiShu(c.FeiShu, fmt.Sprintf(">> %s/%s >> %s",
				c.Name, c.Mode, string(msg.JSON())))
		}
	})

	// or 
	//goio.Init(goio.Environment(c.Mode))
	//goio.SetupLogDefault()
	//goio.Setup(c.FeiShu)
}

```

### 发送通知
#### 飞书
```
普通群消息
gomessage.FeiShu(hookUrl, "这是普通的群消息")
```

#### 钉钉
``` 
gomessage.InitDingDing("token","secret")

普通群消息
err := gomessage.DingDing("这是普通的群消息")

@特定人的消息
@对象必须为绑定钉钉的手机号
err := gomessage.DingDing("Lucy, Harvey, 你们的程序挂了", "18578924567", "+13414567890")

@所有人的消息
err := gomessage.DingDing("这是@所有人的消息", "*")

```

#### Telegram电报
``` 
gomessage.InitTelegram("token",false)

//chatId 个人ID或群组ID text 消息内容
err := gomessage.TelegramTo(chatId, "text")
```

### go-event 基于 chan

```
观察者模式 事件中心
```
```golang
//使用方法
package main

import (
	goevent "github.com/gif-gif/go.io/go-event"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init()
	event := goevent.New()
	event.Subscribe("test", func(msg goevent.Message) {
		golog.WithTag("goevent").Info(msg)
	})
	event.Publish("test", "test")
	time.Sleep(time.Duration(1) * time.Second)
}

```

### GoHttp 模块
```go
package main

import (
	"context"
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/go-http"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	type httpRequest struct {
		Email string `json:"email"`
	}

	req := &gohttp.Request{
		Url: "/main",
		Urls: []string{ //Retry urls
			"/main1",
			"/main2",
			"/main3",
		},
		QueryParams: map[string]string{"name": "jk"},
		Timeout:     time.Second * 2,
		Body: &httpRequest{
			Email: "test@gmail.com",
		},
	}

	gh := &gohttp.GoHttp[gohttp.Response]{
		Request: req,
		BaseUrl: "http://localhost",
		Headers: map[string]string{
			"User-Agent": "github.com/gif-gif/go.io",
		},
	}

	rst, err := gh.HttpPostJson(context.Background())
	if err != nil {
		golog.WithTag("http").Error(err.Error())
	} else {
		fmt.Println(rst)
	}
}

```

### 验证码(支持分布式验证,基于redis)
```go
config := goredis.Config{
    Name:     "gocaptcha",
    Addr:     "127.0.0.1:6379",
    Password: "",
    DB:       0,
    Prefix:   "gocaptcha",
    AutoPing: true,
}

err := gocaptcha.Init(gocaptcha.Config{
    RedisConfig: &config,
})

if err != nil {
    golog.Error(err.Error())
    return
}

data, err := gocaptcha.Default().DigitCaptcha(0, 0, 0)
golog.WithTag("data").Info(data)
```

### gojob 模块

```go
func simpleUseGoJob() {
	n := 0
	cron, err := gojob.New()
	if err != nil {
		golog.WithTag("gojob").Error(err)
	}
	defer cron.Stop()
	cron.Start()

	job, err := cron.SecondX(nil, 1, func(nn int) error {
		golog.WithTag("gojob").Info("testing->" + gconv.String(nn))
		return nil
	}, n)

	if err != nil {
		golog.WithTag("gojob").Error(err)
	} else {
		golog.WithTag("gojob").Info("job.ID:" + job.ID().String())
	}

	time.Sleep(time.Second * 500)
	golog.InfoF("end of gojob")
}
```

### gopool 模块
```go
func testDynamicSize() {
	gp := gopool.NewDynamicSizePool(100, 10)
	defer gp.StopAndWait()

	cron, _ := gojob.New()
	defer cron.Stop()
	cron.Start()
	cron.SecondX(nil, 1, func() {
		gp.PrintPoolStats()
	})

	for i := 0; i < 1000; i++ {
		n := i
		gp.Submit(func() {
			fmt.Printf("Running task #%d\n", n)
			time.Sleep(1 * time.Second)
		})
	}

	golog.InfoF("end of Submit")
}
```
### godb 模块
- mysql
- sqlite3
- clickhouse
- postgresql
- sqlserver
- tidb
```go
func mysqlTest() {
    err := gogorm.Init(gogorm.Config{
     DataSource: "root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Localb",
    })
    if err != nil {
        golog.WithTag("godb").Error(err.Error())
        return
    }
    db := gogorm.Default().DB
	
	err = db.AutoMigrate(&Product{})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}

	// Create
	insertProduct := &Product{Code: "D42", Price: 100}
	db.Create(insertProduct)
	fmt.Println(insertProduct.ID)
	// Read
	var product Product
	tx := db.First(&product, 1) // find product with integer primary key
	if tx.Error != nil {
		fmt.Println("not found first ", tx.Error.Error())
	}
	db.First(&product, "code = ?", "D42")
	// Delete - delete product
	db.Delete(&product, 1)

}
```
- 其他各个模块功能在go-[模块]目录中readme或者testCase中找到使用方法

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
- https://github.com/elastic/go-elasticsearch
- https://github.com/mongodb/mongo-go-driver
- https://github.com/alitto/pond
- https://github.com/tucnak/telebot
- https://github.com/PaulSonOfLars/gotgbot
- https://github.com/mojocn/base64Captcha
- https://github.com/wenlng/go-captcha
- https://github.com/mochi-mqtt/server
- https://github.com/eclipse/paho.mqtt.golang
- https://github.com/apache/rocketmq-client-go
- https://github.com/hibiken/asynq
- https://github.com/xuri/excelize
- https://github.com/samber/lo
- https://github.com/shopspring/decimal
