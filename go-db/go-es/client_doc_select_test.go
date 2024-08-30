package goes

import (
	"log"
	"reflect"
	"testing"
	"time"
)

var server = "http://122.228.113.231:9200"

type TTest struct {
	Id        int `json:"id"`
	UpdatedAt int `json:"updated_at"`
}

type TTT struct {
	Version   string `json:"@version"`
	Id        int    `json:"id"`
	UpdatedAt int    `json:"updated_at"`
	Event     struct {
		Original TTest `json:"original"`
	} `json:"event"`
	Timestamp time.Time `json:"@timestamp"`
}

func TestClient_DocGet(t *testing.T) {
	conf := Config{
		Addr:      server,
		User:      "elastic",
		Password:  "elastic",
		EnableLog: true,
	}
	if err := Init(conf); err != nil {
		log.Println(err.Error())
		return
	}

	index := "logstash-2024.08.29"

	type User struct {
		Name string `json:"name"`
	}

	var u TTT

	if _, err := Default().DocGet(index, "p103nZEBi67xN4wCq_aA", &u); err != nil {
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

	if _, err := Default().DocMultiGet(index, []string{"1001", "1002", "1003"}, &us); err != nil {
		log.Println(err)
		return
	}

	log.Println(us)
}

func TestClient_DocSearch(t *testing.T) {
	conf := Config{
		Addr:      server,
		User:      "elastic",
		Password:  "elastic",
		EnableLog: true,
	}
	if err := Init(conf); err != nil {
		log.Println(err.Error())
		return
	}

	index := "logstash-2024.08.29"

	type User struct {
		Name string `json:"name"`
	}

	query := Default().DocSearchMatchQuery("updated_at", 0)
	p := &Pagination{
		Sort:   "id",
		Order:  true,
		Offset: 0,
		Limit:  10,
	}

	resp, err := Default().DocSearch(index, query, p)
	if err != nil {
		log.Println(err)
		return
	}

	for _, i := range resp.Each(reflect.TypeOf(TTest{})) {
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

	query := Default().DocSearchBoolQuery(BoolQueryFilter,
		Default().DocSearchMatchQuery("name", "goio"),
		Default().DocSearchMatchQuery("name", "noname_1002"))
	p := &Pagination{
		Sort:   "name",
		Order:  false,
		Offset: 0,
		Limit:  10,
	}

	resp, err := Default().DocSearch(index, query, p)
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
