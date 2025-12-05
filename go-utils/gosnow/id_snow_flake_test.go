package gosnow

import (
	"fmt"
	"testing"

	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gogf/gf/util/gconv"
)

// Name     string `yaml:"Name" json:"name,optional"`
// Addr     string `yaml:"Addr" json:"addr,optional"`
// Password string `yaml:"Password" json:"password,optional"`
// DB       int    `yaml:"DB" json:"db,optional"`
// Prefix   string `yaml:"Prefix" json:"prefix,optional"`
func TestGenId(t *testing.T) {
	goredis.Init(goredis.Config{
		Name:     "test",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "snowFake",
	})
	GenIdInit(&SnowFlakeId{WorkerId: 2})
	for i := 0; i < 50; i++ {
		gid := GenIdStr()
		b := goredis.Default().HExists("test", gconv.String(gid)).Val()
		if b {
			golog.WithTag("snowFlake").Info("my god")
		} else {
			goredis.Default().HSet("test", gconv.String(gid), "1")
		}

		fmt.Println(gid)
	}
	golog.WithTag("snowFlake").Info("The End")
}
