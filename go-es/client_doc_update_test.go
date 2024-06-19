package goo_es

import (
	"fmt"
	"log"
	"testing"
)

func TestClient_DocUpdate(t *testing.T) {
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

	fmt.Println(__client.DocUpdate(index, "1001", map[string]interface{}{"name": "noname_1001"}))
}

func TestClient_DocUpset(t *testing.T) {
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

	fmt.Println(__client.DocUpset(index, "1004", map[string]interface{}{"name": "noname_1004"}))
}

func TestClient_DocBatchUpdate(t *testing.T) {
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

	var (
		ids  []string
		docs []interface{}
	)

	for i := 1; i < 10; i++ {
		ids = append(ids, fmt.Sprintf("%d", 1000+i))
		docs = append(docs, map[string]interface{}{
			"name": fmt.Sprintf("noname_%d", 1000+i),
		})
	}

	fmt.Println(__client.DocBatchUpset(index, ids, docs))
}

func TestClient_DocBatchUpset(t *testing.T) {
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

	var (
		ids  []string
		docs []interface{}
	)

	for i := 1; i < 10; i++ {
		ids = append(ids, fmt.Sprintf("%d", 1000+i))
		docs = append(docs, map[string]interface{}{
			"name": fmt.Sprintf("hnatao_%d", 1000+i),
		})
	}

	fmt.Println(__client.DocBatchUpset(index, ids, docs))
}
