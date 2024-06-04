package goio

import golog "github.com/gif-gif/go.io/go-log"

// 当前运行环境
var Env Environment

func Init(env Environment) {
	Env = env
	golog.SetAdapter(golog.NewFileAdapter()) //默认当前工程目录logs/date.log
}
