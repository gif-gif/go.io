package goes

import (
	"errors"
	"github.com/olivere/elastic/v7"
)

// 文档 - 添加
// index: 索引
// id: 唯一标识
// body: json格式的数据
func (cli *GoEs) DocCreate(index, id string, body interface{}) (*elastic.IndexResponse, error) {
	return cli.cli.Index().Type("_doc").
		Index(index).
		OpType("create").
		Id(id).
		BodyJson(body).
		Refresh("true").
		Do(cli.ctx)
}

// 文档 - 批量添加
func (cli *GoEs) DocBatchCreate(index string, data map[string]interface{}) (resp *elastic.BulkResponse, err error) {
	bs := cli.cli.Bulk().Index(index).Refresh("true")
	for id, doc := range data {
		bs.Add(elastic.NewBulkIndexRequest().Id(id).Doc(doc)).Index(index)
	}

	resp, err = bs.Do(cli.ctx)
	if err != nil {
		return
	}

	if l := len(resp.Failed()); l > 0 {
		err = errors.New(resp.Failed()[0].Error.Reason)
	}
	return
}
