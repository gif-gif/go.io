package gokafka

import goredis "github.com/gif-gif/go.io/go-db/go-redis"

const (
	FocusName = "focus"
	RedisName = "redis"
)

type Option struct {
	Name  string
	Value interface{}
}

// 是否强制
func FocusOption() Option {
	return Option{Name: FocusName, Value: true}
}

// redis 对象
func RedisOption(cli *goredis.GoRedis) Option {
	return Option{Name: RedisName, Value: cli}
}
