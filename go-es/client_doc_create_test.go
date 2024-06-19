package goo_es

import (
	"fmt"
	"log"
	"testing"
)

func TestClient_DocCreate(t *testing.T) {
	conf := Config{
		Addr:     "http://192.168.1.100:9200",
		User:     "elastic",
		Password: "123456",
	}
	if err := Init(conf); err != nil {
		log.Println(err.Error())
		return
	}

	index := "test_202209"

	resp, err := __client.DocCreate(index, "1001", map[string]interface{}{
		"name": "hnatao",
	})
	fmt.Println(resp, err)
}

func TestClient_DocBatchCreate(t *testing.T) {
	conf := Config{
		Addr:     "http://192.168.1.100:9200",
		User:     "elastic",
		Password: "123456",
	}
	if err := Init(conf); err != nil {
		log.Println(err.Error())
		return
	}

	index := "test_202209"

	data := map[string]interface{}{
		"1002": map[string]interface{}{
			"name": "noname_1002",
		},
		"1003": map[string]interface{}{
			"name": "noname_1003",
		},
	}
	resp, err := __client.DocBatchCreate(index, data)
	fmt.Println(resp, err)
}
