package goio

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
