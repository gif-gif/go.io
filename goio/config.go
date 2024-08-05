package goio

import (
	"encoding/json"
	golog "github.com/gif-gif/go.io/go-log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadYamlConfig(yamlFile string, conf interface{}) (err error) {
	var buf []byte

	buf, err = ioutil.ReadFile(yamlFile)
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
	var buf []byte

	buf, err = ioutil.ReadFile(jsonFile)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	if err = json.Unmarshal(buf, conf); err != nil {
		golog.Error(err.Error())
	}
	return
}
