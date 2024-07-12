package goio

// 当前运行环境
var Env Environment //string `json:",default=pro,options=dev|test|rt|pre|pro"`

func Init(env Environment) {
	Env = env
	//golog.SetAdapter(golog.NewFileAdapter()) //默认当前工程目录logs/date.log
}
