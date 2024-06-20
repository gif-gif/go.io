package goes

import (
	"github.com/olivere/elastic"
)

// 文档 - 删除 - 根据ID
func (cli *client) DocDelete(index, id string) (*elastic.IndexResponse, error) {
	return cli.cli.Index().Index(index).Id(id).Refresh("true").Do(cli.ctx)
}

// 文档 - 删除 - 根据条件
// 即时没有符合条件的文档，也不会报404错误
func (cli *client) DocDeleteBy(index string, query elastic.Query) (total int64, err error) {
	var resp *elastic.BulkIndexByScrollResponse
	resp, err = cli.cli.DeleteByQuery(index).Query(query).Refresh("true").Do(cli.ctx)
	if err != nil {
		return
	}
	total = resp.Deleted
	return
}

// 文档 - 批量删除
func (cli *client) DocBatchDelete(index string, ids []string) (total int64, err error) {
	if l := len(ids); l == 0 {
		return
	}

	bs := cli.cli.Bulk().Index(index).Refresh("true")
	for _, id := range ids {
		bs.Add(elastic.NewBulkDeleteRequest().Id(id))
	}

	var resp *elastic.BulkResponse
	resp, err = bs.Do(cli.ctx)
	if err != nil {
		return
	}

	total = int64(len(resp.Succeeded()))
	return
}
