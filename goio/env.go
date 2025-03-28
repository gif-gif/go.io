package goio

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/go-log/adapters"
	gomessage "github.com/gif-gif/go.io/go-message"
)

type Environment string

// 库运行模式 dev|test|rt|pre|pro go-zero
const (
	PRODUCTION  Environment = "pro"
	PRE         Environment = "pre"
	RT          Environment = "rt"
	TEST        Environment = "test"
	DEVELOPMENT Environment = "dev"
)

var (
	envTags = map[Environment]string{
		PRODUCTION:  "prod",
		PRE:         "pre",
		RT:          "rt",
		TEST:        "test",
		DEVELOPMENT: "dev",
	}
)

func (env Environment) String() string {
	return string(env)
}

func (env Environment) Tag() string {
	return envTags[env]
}

// 日志级别一键设置， 测试环境和生产环境发送飞书消息
//
// 也可以不用这个方法，项目中可以自定义， golog.WithHook 多次调用不会覆盖，每次调用都会生效
func Setup(feishu string) {
	// 发送飞书消息
	if feishu == "" {
		return
	}
	golog.WithHook(func(msg *golog.Message) {
		if msg.Level <= golog.WARN {
			return
		}
		if Env == TEST || Env == PRODUCTION {
			gomessage.FeiShu(feishu, fmt.Sprintf(">> %s/%s >> %s", Name, Env.String(), string(msg.JSON())))
		}
	})
}

func SetupLogDefault() {
	if Env == TEST || Env == PRODUCTION {
		golog.SetAdapter(adapters.NewFileAdapter())
	}
}

func IsTest() bool {
	return Env == TEST
}

func IsPro() bool {
	return Env == PRODUCTION
}

func IsDev() bool {
	return Env == DEVELOPMENT
}
