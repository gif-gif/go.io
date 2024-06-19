package goo_es

import (
	"log"
	"reflect"
	"testing"
)

func TestClient_DocGet(t *testing.T) {
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

	type User struct {
		Name string `json:"name"`
	}

	var u User

	if _, err := __client.DocGet(index, "1001", &u); err != nil {
		log.Println(err)
		return
	}

	log.Println(u)
}

func TestClient_DocMultiGet(t *testing.T) {
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

	type User struct {
		Name string `json:"name"`
	}

	var us []User

	if _, err := __client.DocMultiGet(index, []string{"1001", "1002", "1003"}, &us); err != nil {
		log.Println(err)
		return
	}

	log.Println(us)
}

func TestClient_DocSearch(t *testing.T) {
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

	type User struct {
		Name string `json:"name"`
	}

	query := __client.DocSearchMatchQuery("name", "hnatao")
	p := &Pagination{
		Sort:   "name",
		Order:  true,
		Offset: 0,
		Limit:  10,
	}

	resp, err := __client.DocSearch(index, query, p)
	if err != nil {
		log.Println(err)
		return
	}

	for _, i := range resp.Each(reflect.TypeOf(User{})) {
		log.Println(i)
	}
	log.Println(resp.Hits.TotalHits)
	log.Println(resp.Hits.Hits)
}

func TestClient_DocSearch2(t *testing.T) {
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

	type User struct {
		Name string `json:"name"`
	}

	query := __client.DocSearchBoolQuery(BoolQueryFilter,
		__client.DocSearchMatchQuery("name", "hnatao"),
		__client.DocSearchMatchQuery("name", "noname_1002"))
	p := &Pagination{
		Sort:   "name",
		Order:  false,
		Offset: 0,
		Limit:  10,
	}

	resp, err := __client.DocSearch(index, query, p)
	if err != nil {
		log.Println(err)
		return
	}
	for _, i := range resp.Each(reflect.TypeOf(User{})) {
		log.Println(i)
	}

	log.Println(resp.Hits.TotalHits)
	log.Println(resp.Hits.Hits)
}
