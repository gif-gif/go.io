package goes

import (
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
)

// 文档 - 修改
func (cli *client) DocUpdate(index, id string, body interface{}) (response *elastic.UpdateResponse, err error) {
	return cli.cli.Update().
		Index(index).
		Id(id).
		Doc(body).
		Refresh("true").
		Do(cli.ctx)
}

// 文档 - 修改 - 不存在就插入
func (cli *client) DocUpset(index, id string, body interface{}) (response *elastic.UpdateResponse, err error) {
	return cli.cli.Update().
		Index(index).
		Id(id).
		Doc(body).
		Upsert(body).
		Refresh("true").
		Do(cli.ctx)
}

// 文档 - 修改 - 根据条件
func (cli *client) DocUpdateBy(index string, query elastic.Query, script *elastic.Script) (total int64, err error) {
	var resp *elastic.BulkIndexByScrollResponse
	resp, err = cli.cli.UpdateByQuery(index).Query(query).Script(script).Refresh("true").Do(cli.ctx)
	if err != nil {
		return
	}
	total = resp.Updated
	return
}

// 文档 - 批量修改
func (cli *client) DocBatchUpdate(index string, ids []string, docs []interface{}) (err error) {
	bs := cli.cli.Bulk().Index(index).Refresh("true")
	for i := range ids {
		bs.Add(elastic.NewBulkUpdateRequest().Id(ids[i]).Doc(docs[i]))
	}
	var resp *elastic.BulkResponse
	resp, err = bs.Do(cli.ctx)
	if l := len(resp.Failed()); l > 0 {
		err = errors.New(resp.Failed()[0].Error.Reason)
		return
	}
	return
}

// 文档 - 批量修改 - 不存在就插入
func (cli *client) DocBatchUpset(index string, ids []string, docs []interface{}) (err error) {
	bs := cli.cli.Bulk().Index(index).Refresh("true")
	for i := range ids {
		bs.Add(elastic.NewBulkUpdateRequest().Id(ids[i]).Doc(docs[i]).Upsert(docs[i]))
	}
	var resp *elastic.BulkResponse
	resp, err = bs.Do(cli.ctx)
	if l := len(resp.Failed()); l > 0 {
		err = errors.New(resp.Failed()[0].Error.Reason)
		return
	}
	return
}

// -----------------------------------------------------------------------
// 脚本更新示例
// 参考文档：https://dablelv.blog.csdn.net/article/details/121396060
// -----------------------------------------------------------------------

// 1. 数组删除元素
func (cli *client) DocArrayDelField(field string, value interface{}) *elastic.Script {
	scriptStr := fmt.Sprintf(
		`ctx._source.%s.remove(ctx._source.%s.indexOf(params.%s))`,
		field, field, field,
	)
	return elastic.NewScript(scriptStr).Params(
		map[string]interface{}{
			field: value,
		},
	)
}

// 2. 数组删除多个元素
func (cli *client) DocArrayDelFields(field string, value []interface{}) *elastic.Script {
	scriptStr := fmt.Sprintf(
		`for (int i = 0; i < params.%s.length; i++) {
					if (ctx._source.%s.contains(params.%s[i])) { 	
						ctx._source.%s.remove(ctx._source.%s.indexOf(params.%s[i]))
					}
				}`,
		field, field, field, field, field, field,
	)
	return elastic.NewScript(scriptStr).Params(
		map[string]interface{}{
			field: value,
		},
	)
}

// 3. 数组追加元素
func (cli *client) DocArrayAppendValue(field string, value interface{}) *elastic.Script {
	scriptStr := fmt.Sprintf(
		`ctx._source.%s.add(params.%s)`,
		field, field,
	)
	return elastic.NewScript(scriptStr).Params(
		map[string]interface{}{
			field: value,
		},
	)
}

// 4. 数组追加多个元素
func (cli *client) DocArrayAppendValues(field string, value []interface{}) *elastic.Script {
	scriptStr := fmt.Sprintf(
		`ctx._source.%s.addAll(params.%s)`,
		field, field,
	)
	return elastic.NewScript(scriptStr).Params(
		map[string]interface{}{
			field: value,
		},
	)
}

// 5. 数组修改元素
func (cli *client) DocArrayUpdateValue(field string, old, new interface{}) *elastic.Script {
	scriptStr := fmt.Sprintf(
		`ctx._source.%s[ctx._source.%s.indexOf(params.old)]=params.new`,
		field, field,
	)
	return elastic.NewScript(scriptStr).Params(
		map[string]interface{}{
			"old": old,
			"new": new,
		},
	)
}
