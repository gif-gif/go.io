package goes

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

	//if err := Default().IndexCreate(index, ``); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//fmt.Println(Default().IndexGet(index))
	//fmt.Println(Default().IndexMapping(index))
	//fmt.Println(Default().IndexSettings(index))

	fmt.Println(Default().IndexUpdateMapping(index, `{"properties": {"name": {"type":"text", "fielddata": true}}}`))

	//fmt.Println(Default().IndexExists(index))

	//fmt.Println(Default().IndexNames())

	//fmt.Println(Default().IndexAlias(index, index+"_aa"))
	//fmt.Println(Default().IndexAliasRemove(index, index+"_aa"))

	//fmt.Println(Default().Client().ElasticsearchVersion(conf.Addr))
}
