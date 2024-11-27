package goio

// 当前运行环境
var Env Environment = "dev" //string `json:",default=pro,options=dev|test|rt|pre|pro"`
var Name string

func Init(env Environment, name ...string) {
	Env = env
	if l := len(name); l > 0 {
		for i, n := range name {
			if n != "" {
				if i == 0 {
					Name = n
				} else {
					Name += "_" + n
				}
				break
			}
		}
	} else {
		Name = "default"
	}
	//golog.SetAdapter(golog.NewFileAdapter()) //默认当前工程目录logs/date.log
}
