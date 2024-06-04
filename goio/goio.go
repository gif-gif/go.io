package goio

// 当前运行环境
var Env Environment

func Init(env Environment) {
	Env = env
}
