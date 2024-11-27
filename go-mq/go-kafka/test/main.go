package main

import (
	"encoding/json"
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
	goutils "github.com/gif-gif/go.io/go-utils"
	"gorm.io/gorm"
	"testing"
	"time"
)

var (
	topic   = "test-01"
	groupId = "test-01"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type Account struct {
	Id        int64 `json:"id"`
	UpdatedAt int64 `json:"updated_at"`
}

type TestMessage struct {
	Id      int    `json:"id"`
	TraceId string `json:"trace_id"`
}

func (t *TestMessage) Topic() string {
	return topic
}

func (t *TestMessage) Key() string {
	return fmt.Sprintf("%s:%d", topic, t.Id)
}

func (t *TestMessage) Headers() map[string]string {
	return map[string]string{
		"source": "my-test",
	}
}

func (t *TestMessage) Serialize() []byte {
	t.TraceId = goutils.UUID()
	b, _ := json.Marshal(t)
	return b
}

func (t *TestMessage) Deserialize(b []byte) {
	if err := json.Unmarshal(b, t); err != nil {
		fmt.Println(err)
	}
}

func TestProducer(t *testing.T) {
	redis, _ := goredis.New(goredis.Config{
		Addr:     "redis.in:20063",
		Password: "",
		DB:       0,
	})
	gokafka.Init(gokafka.Config{
		Addrs:        []string{"122.228.113.231:30094"}, //122.228.113.231
		User:         "admin",
		Password:     "b36da6b4eb0f3",
		Timeout:      10,
		OffsetNewest: false,
	}, gokafka.RedisOption(redis))

	for i := 0; i < 20; i++ {
		gokafka.Producer().SendMessage(&TestMessage{Id: 200 + i})
	}

	time.Sleep(3 * time.Second)
}

func TestConsumer(t *testing.T) {
	gokafka.Init(gokafka.Config{
		User:     "admin",
		Password: "",
		Addrs:    []string{"kafka.in:20092"},
		RedisConfig: goredis.Config{
			Addr:     "redis.in:20063",
			Password: "",
			DB:       0,
		},
	})

	gokafka.Consumer().ConsumeGroup(groupId, []string{topic}, func(ctx *gocontext.Context, msg *gokafka.ConsumerMessage, consumerErr *gokafka.ConsumerError) error {
		time.Sleep(5 * time.Second)
		return nil
	})
}
