package goio

type Env string

// 库运行模式 dev|test|rt|pre|pro go-zero
const (
	PRODUCTION  Env = "pro"
	PRE         Env = "pre"
	RT          Env = "rt"
	TEST        Env = "test"
	DEVELOPMENT Env = "dev"
)

var (
	envTags = map[Env]string{
		PRODUCTION:  "prod",
		PRE:         "pre",
		RT:          "rt",
		TEST:        "test",
		DEVELOPMENT: "dev",
	}
)

func (env Env) String() string {
	return string(env)
}

func (env Env) Tag() string {
	return envTags[env]
}
