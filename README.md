# go.io
- Golang Development Framework Continuously Developing and Updating

## 设计目标
- goio 提供了常用库封装，支持必要的简洁使用功能，在其之上可以进二次开发，以提供更好的代码维护；
- 以跨平台跨项目为首要原则，以减少二次开发的成本；
- 各个模块逻辑保持唯一不重复,一定程度上项目以第三方库解偶

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
}

```

### 发送通知
飞书
```
普通群消息
gomessage.FeiShu(hookUrl, "这是普通的群消息")
```

钉钉
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
	"fmt"
	gohttpx "github.com/gif-gif/go.io/go-http/go-httpex"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)

	req := gohttpx.Request{
		Url: "http://localhost:100",
		Urls: []string{
			"http://localhost:200",
			"http://localhost:300",
			"http://localhost:400",
		},
		QueryParams: map[string]string{"name": "jk"},
		Timeout:     time.Second * 2,
	}
	type httpRequest struct {
		Email string `json:"email"`
	}

	req.Body = &httpRequest{
		Email: "test@gmail.com",
	}

	res := &gohttpx.Response{}
	err := gohttpx.HttpPostJson[gohttpx.Response](&req, res)
	if err != nil {
		golog.ErrorF("Error: %+v\n", err)
	} else {
		fmt.Println(res)
	}

	time.Sleep(10 * time.Second)
}

```
### 验证码(支持分布式验证,基于redis)
```
config := goredis.Config{
    Name:     "gocaptcha-goredis",
    Addr:     "127.0.0.1:6379",
    Password: "",
    DB:       0,
    Prefix:   "gocaptcha",
    AutoPing: true,
}

a, err := gocaptcha.NewRedis(config)
if err != nil {
    golog.Error(err.Error())
    return
}
goutils.InitCaptcha(a) //不初始化默认 memorycache
goutils.CaptchaGet(200,40)
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
	db, err := gogorm.InitMysql("root:223238@tcp(127.0.0.1:33060)/gromdb?charset=utf8mb4&parseTime=True&loc=Local", godb.GoDbConfig{})
	if err != nil {
		golog.WithTag("godb").Error(err.Error())
		return
	}
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

## 杂项

|Status|Value|
|:----|:---:|
|Stars|[![Stars](https://img.shields.io/github/stars/gif-gif/go.io)](https://github.com/gif-gif/go.io)
|Forks|[![Forks](https://img.shields.io/github/forks/gif-gif/go.io)](https://github.com/gif-gif/go.io)
|License|[![License](https://img.shields.io/github/license/gif-gif/go.io)](https://github.com/gif-gif/go.io)
|Issues|[![Issues](https://img.shields.io/github/issues/gif-gif/go.io)](https://github.com/gif-gif/go.io)
|Release Downloads|[![Downloads](https://img.shields.io/github/downloads/gif-gif/go.io/total.svg)](https://github.com/gif-gif/go.io/releases)
