package test

import (
	"encoding/json"
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	goetcd "github.com/gif-gif/go.io/go-etcd"
	etcd_configurator "github.com/gif-gif/go.io/go-zero/etcd-configurator"
	"log"
	"testing"
)

// 配置结构定义
type TestSt struct {
	Name string `json:"name"`
}

func TestEtcdConfigListener(t *testing.T) {
	etcd_configurator.NewConfigCenter[TestSt]("config-test", goetcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}, func(t TestSt) {
		println(t.Name)
	})

	<-gocontext.Cancel().Done()
}

func TestEtcdSaveConfig(t *testing.T) {
	goetcd.Init(goetcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		//Username:  "root",
		//Password:  "123456",
	})
	goetcd.Del("config-test")

	data := &TestSt{
		Name: "Test112",
	}
	str, _ := json.Marshal(data)
	if _, err := goetcd.Set("config-test", string(str)); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Setting config:", goetcd.GetString("config-test"))
	//Del("config-test")
}
