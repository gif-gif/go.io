package goeso

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type GoEs struct {
	es  *elasticsearch.TypedClient
	esc *elasticsearch.Client
	ctx context.Context
}

func New(conf elasticsearch.Config) (cli *GoEs, err error) {
	cli = &GoEs{
		ctx: context.Background(),
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		return nil, err
	}

	cli.esc = es
	err = cli.newType(conf)
	if err != nil {
		return nil, err
	}
	//
	//goutils.AsyncFunc(func() {
	//	//cli.es.Ping = func() *ping.Ping {
	//	//	return nil
	//	//}
	//})

	return cli, nil
}

func (cli *GoEs) newType(conf elasticsearch.Config) (err error) {
	es, err := elasticsearch.NewTypedClient(conf)
	if err != nil {
		return err
	}

	cli.es = es
	return nil
}

func (cli *GoEs) TypeClient() *elasticsearch.TypedClient {
	return cli.es
}

func (cli *GoEs) Client() *elasticsearch.Client {
	return cli.esc
}

// 文档 - 添加
// index: 索引
// docId: 唯一标识
// body: json or struct
func (cli *GoEs) DocCreate(index, docId string, document interface{}) (*esapi.Response, error) {
	if value, ok := document.(string); ok {
		return cli.esc.Create(index, docId, strings.NewReader(value))
	} else {
		v, _ := json.Marshal(document)
		return cli.esc.Create(index, docId, strings.NewReader(string(v)))
	}
}

// Updating documents
// client.Update("my_index", "id", strings.NewReader(`{doc: { language: "Go" }}`))
func (cli *GoEs) DocUpdate(index, docId string, document interface{}) (*esapi.Response, error) {
	if value, ok := document.(string); ok {
		return cli.esc.Update(index, docId, strings.NewReader(value))
	} else {
		v, _ := json.Marshal(document)
		return cli.esc.Update(index, docId, strings.NewReader(string(v)))
	}
}

// 文档 - 批量添加
func (cli *GoEs) DocBatchCreate(index string, data map[string]interface{}) (err error) {
	indexer, _ := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: cli.esc,
		Index:  index,
	})

	for id, doc := range data {
		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: id,
				Body:       strings.NewReader(gconv.String(doc)),
			})
		if err != nil {
			return err
		}
	}

	err = indexer.Close(cli.ctx)
	if err != nil {
		return err
	}

	return nil
}

// 删除 - 根据IDS - Deleting an index
func (cli *GoEs) DeleteIndex(indexes []string) (*esapi.Response, error) {
	return cli.esc.Indices.Delete(indexes)
}

// 文档 - 删除 - 根据ID - Deleting documents
func (cli *GoEs) DocDelete(index, id string) (*esapi.Response, error) {
	return cli.esc.Delete(index, id)
}

// 文档 - 删除 - 根据条件
// 即时没有符合条件的文档，也不会报404错误
func (cli *GoEs) DocDeleteBy(index string, query string) (*esapi.Response, error) {
	return cli.esc.DeleteByQuery([]string{index}, strings.NewReader(query))
}

// 文档 - 批量删除
func (cli *GoEs) DocBatchDelete(index string, ids []string) (total int64, err error) {
	if l := len(ids); l == 0 {
		return
	}
	total = 0
	for _, id := range ids {
		_, err = cli.DocDelete(index, id)
		if err == nil {
			total++
		} else {
			return total, err
		}
	}

	return total, nil
}
