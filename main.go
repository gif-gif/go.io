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

// 初始化实例
func main() {
	c := Config{}
	goio.Init(goio.Environment(c.Mode))
	//golog.SetAdapter(golog.NewFileAdapter()) //当前工程目录logs/date.log, 可通过这个设置改变日志输出目录
	//不设置 golog.SetAdapter 默认控制台输出
	golog.WithHook(func(msg *golog.Message) {
		if msg.Level > golog.ERROR { //致命错误以上发通知
			gomessage.FeiShu(c.FeiShu, fmt.Sprintf(">> %s/%s >> %s",
				c.Name, c.Mode, string(msg.JSON())))
		}
	})

	golog.WithTag("goio").ErrorF("Starting tests")
}
