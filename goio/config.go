package goio

import (
	"encoding/json"
	gomongo "github.com/gif-gif/go.io/go-db/go-mongo"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/go-utils/prometheusx"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Env Environment `yaml:"env"`

	Server struct {
		Addr string `yaml:"addr"`
		Name string `yaml:"name"`
	} `yaml:"server"`

	Prometheus prometheusx.Config `yaml:"prometheus"`
	MongoDB    gomongo.Config     `yaml:"mongodb,omitempty"`
	//Mysql       gogorm.Config  `yaml:"mysql,omitempty"`
	//Postgres    gogorm.Config  `yaml:"postgres,omitempty"`
	//Sqlite      gogorm.Config  `yaml:"sqlite,omitempty"`
	//Clickhouse1 gogorm.Config  `yaml:"clickhouse1,omitempty"`
	//Redis       goredis.Config      `yaml:"redis,omitempty"`
	//Kafka       gokafka.Config      `yaml:"kafka,omitempty"`
	//Clickhouse  goclickhouse.Config `yaml:"clickhouse,omitempty"`
	//Es          goes.Config         `yaml:"es,omitempty"`
	////EsIndex EsIndex          `yaml:"es_index"`
	//
	//Etcd goetcd.Config `yaml:"etcd"`

	FeiShu string `yaml:"feishu"`
}

func LoadYamlConfig(yamlFile string, conf interface{}) (err error) {
	if yamlFile == "" {
		yamlFile = ".yaml"
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
