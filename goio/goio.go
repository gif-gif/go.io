package goio

import (
	"context"
	gocontext "github.com/gif-gif/go.io/go-context"
)

// 当前运行环境
var Env Environment = "dev" //string `json:",default=pro,options=dev|test|rt|pre|pro"`

func Init(env Environment) {
	Env = env
	//golog.SetAdapter(golog.NewFileAdapter()) //默认当前工程目录logs/date.log
}

func Context() context.Context {
	return gocontext.Cancel()
}
