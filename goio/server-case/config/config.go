package conf

import (
	"encoding/json"
	goclickhouse "github.com/gif-gif/go.io/go-db/go-clickhouse"
	goes "github.com/gif-gif/go.io/go-db/go-es"
	gomongo "github.com/gif-gif/go.io/go-db/go-mongo"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	"github.com/gif-gif/go.io/go-db/gogorm"
	goetcd "github.com/gif-gif/go.io/go-etcd"
	golog "github.com/gif-gif/go.io/go-log"
	gokafka "github.com/gif-gif/go.io/go-mq/go-kafka"
	"github.com/gif-gif/go.io/go-utils/prometheusx"
	"github.com/gif-gif/go.io/goio"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Env goio.Environment `yaml:"env"`

	Server struct {
		Addr string `yaml:"addr"`
		Name string `yaml:"name"`
	} `yaml:"server"`

	Prometheus  prometheusx.Config  `yaml:"prometheus"`
	MongoDB     gomongo.Config      `yaml:"mongodb"`
	Mysql       gogorm.Config       `yaml:"mysql"`
	Postgres    gogorm.Config       `yaml:"postgres"`
	Sqlite      gogorm.Config       `yaml:"sqlite"`
	Clickhouse1 gogorm.Config       `yaml:"clickhouse1"`
	Redis       goredis.Config      `yaml:"redis"`
	Kafka       gokafka.Config      `yaml:"kafka"`
	Clickhouse  goclickhouse.Config `yaml:"clickhouse"`
	Es          goes.Config         `yaml:"es"`
	////EsIndex EsIndex          `yaml:"es_index"`

	Etcd goetcd.Config `yaml:"etcd"`

	FeiShu string `yaml:"feishu"`
}

func LoadYamlConfig(yamlFile string, conf interface{}) (err error) {
	if yamlFile == "" {
		yamlFile = "api-local.yaml"
	}

	var buf []byte

	buf, err = os.ReadFile(yamlFile)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	if err = yaml.Unmarshal(buf, conf); err != nil {
		golog.Error(err.Error())
	}
	return
}

func LoadJsonConfig(jsonFile string, conf interface{}) (err error) {
	if jsonFile == "" {
		jsonFile = ".json"
	}
	var buf []byte

	buf, err = os.ReadFile(jsonFile)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	if err = json.Unmarshal(buf, conf); err != nil {
		golog.Error(err.Error())
	}
	return
}
