package goo_es

import (
	"fmt"
	"log"
	"testing"
)

func TestClient_IndexGet(t *testing.T) {
	conf := Config{
		Addr:      "http://192.168.1.100:9200",
		User:      "elastic",
		Password:  "123456",
		EnableLog: true,
	}
	if err := Init(conf); err != nil {
		log.Println(err.Error())
		return
	}

	index := "test_202209"

	//if err := __client.IndexCreate(index, ``); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//fmt.Println(__client.IndexGet(index))
	//fmt.Println(__client.IndexMapping(index))
	//fmt.Println(__client.IndexSettings(index))

	fmt.Println(__client.IndexUpdateMapping(index, `{"properties": {"name": {"type":"text", "fielddata": true}}}`))

	//fmt.Println(__client.IndexExists(index))

	//fmt.Println(__client.IndexNames())

	//fmt.Println(__client.IndexAlias(index, index+"_aa"))
	//fmt.Println(__client.IndexAliasRemove(index, index+"_aa"))

	//fmt.Println(__client.Client().ElasticsearchVersion(conf.Addr))
}
