package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	goes "github.com/gif-gif/go.io/go-db/go-esnew"
	golog "github.com/gif-gif/go.io/go-log"
)

type document struct {
	Name string `json:"name"`
}

func main() {
	testSearch()
	golog.Info("finish...")
}

func testSearch() {
	//es, err := goes.New(elasticsearch.Config{
	//	Addresses:     []string{"http://122.228.113.238:9200"},
	//	Username:      "es",
	//	Password:      "123456",
	//	RetryOnStatus: []int{502, 503, 504, 429},
	//})
	//
	//if err != nil {
	//	golog.Error("Error creating elasticsearch:" + err.Error())
	//	return
	//}
	//
	//doc := document{}
	//index := "index_test"
	//query := `{ "query": { "match_all": {} } }`
	//
	//_, err = es.DocSearch(index, types.Query{MatchAll: &types.MatchAllQuery{}}, nil)
	//res, err := es.DocMultiGet(index, []string{"1", "2"})
	//if err != nil {
	//	return
	//}
	//遍历所有结果
	//for _, hit := range res.Hits.Hits {
	//	fmt.Printf("%s\n", hit.Source_)
	//}
	//
	//_, err = es.DocCreate(index, "2", str)
	//if err != nil {
	//	golog.Error("Error DocCreate elasticsearch:" + err.Error())
	//	return
	//}
	//
	//_, err = es.DeleteIndex([]string{"index1"})
	//if err != nil {
	//	golog.Error("Error deleting elasticsearch:" + err.Error())
	//	return
	//}

}

func testCreateDoc() {
	es, err := goes.New(elasticsearch.Config{
		Addresses:     []string{"http://122.228.113.238:9200"},
		Username:      "es",
		Password:      "123456",
		RetryOnStatus: []int{502, 503, 504, 429},
	})

	if err != nil {
		golog.Error("Error creating elasticsearch:" + err.Error())
		return
	}

	doc := document{Name: "test"}
	index := "test"
	_, err = es.DocCreate(index, "2", doc)
	if err != nil {
		golog.Error("Error DocCreate elasticsearch:" + err.Error())
		return
	}

	//_, err = es.DeleteIndex([]string{"index1"})
	//if err != nil {
	//	golog.Error("Error deleting elasticsearch:" + err.Error())
	//	return
	//}

	golog.Info("finish...")
}

func testDocGet() {
	es, err := goes.New(elasticsearch.Config{
		Addresses:     []string{"http://122.228.113.238:9200"},
		Username:      "es",
		Password:      "123456",
		RetryOnStatus: []int{502, 503, 504, 429},
	})

	if err != nil {
		golog.Error("Error creating elasticsearch:" + err.Error())
		return
	}
	type document struct {
		Name string `json:"name"`
	}

	doc := document{}
	index := "index_test"
	_, err = es.DocGet(index, "1", &doc)
	if err != nil {
		golog.Error("Error DocGet elasticsearch:" + err.Error())
		return
	} else {
		golog.Info(doc.Name)
	}
}
