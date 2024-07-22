package goes

import (
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
)

// 文档 - 查询文档 - 根据ID
func (cli *client) DocGet(index, id string, m interface{}) (resp *elastic.GetResult, err error) {
	resp, err = cli.cli.Get().Index(index).Id(id).Do(cli.ctx)
	if err != nil {
		return
	}
	if resp.Error != nil {
		err = errors.New(resp.Error.Reason)
		return
	}
	if m != nil {
		json.Unmarshal(resp.Source, &m)
	}
	return
}

// 文档 - 查询文档 - 根据ID集合
func (cli *client) DocMultiGet(index string, ids []string, m interface{}) (resp *elastic.MgetResponse, err error) {
	ms := cli.cli.MultiGet()
	for _, id := range ids {
		ms.Add(elastic.NewMultiGetItem().Index(index).Id(id))
	}
	resp, err = ms.Do(cli.ctx)
	if err != nil {
		return
	}
	if m != nil {
		var arr []interface{}
		for _, doc := range resp.Docs {
			var _m interface{}
			json.Unmarshal(doc.Source, &_m)
			arr = append(arr, _m)
		}
		b, _ := json.Marshal(&arr)
		json.Unmarshal(b, &m)
	}
	return
}

type Pagination struct {
	Sort   string // 排序字段
	Order  bool   // true=升序 false=降序
	Offset int
	Limit  int
}

// 文档 - 搜索文档 - 根据条件
func (cli *client) DocSearch(index string, query elastic.Query, p *Pagination) (resp *elastic.SearchResult, err error) {
	ss := cli.cli.Search().Index(index).Query(query)

	if p != nil {
		if v := p.Sort; v != "" {
			ss.Sort(v, p.Order)
		}
		if v := p.Offset; v > 0 {
			ss.From(v)
		}
		if v := p.Limit; v > 0 {
			ss.Size(v)
		}
	}

	resp, err = ss.Do(cli.ctx)
	if err != nil {
		return
	}
	if resp.Error != nil {
		err = errors.New(resp.Error.Reason)
	}
	return
}

// 文档 - 搜索条件 - 精确匹配单个字段
func (cli *client) DocSearchTermQuery(field string, value interface{}) *elastic.TermQuery {
	return elastic.NewTermQuery(field, value)
}

// 文档 - 搜索条件 - 精确匹配多个字段
func (cli *client) DocSearchTermsQuery(field string, values []interface{}) *elastic.TermsQuery {
	return elastic.NewTermsQuery(field, values...)
}

// 文档 - 搜索条件 - 匹配查找，单字段搜索（匹配分词结果）
func (cli *client) DocSearchMatchQuery(field string, value interface{}) *elastic.MatchQuery {
	return elastic.NewMatchQuery(field, value)
}

// 文档 - 搜索条件 - 脚本查询
func (cli *client) DocSearchScriptQuery(script *elastic.Script) *elastic.ScriptQuery {
	return elastic.NewScriptQuery(script)
}

// 文档 - 搜索条件 - 范围查找
func (cli *client) DocSearchRangeQuery(field string, from, to interface{}) *elastic.RangeQuery {
	return elastic.NewRangeQuery(field).Gte(from).Lte(to)
}

// 文档 - 搜索条件 - 判断某个字段是否存在
func (cli *client) DocSearchExistsQuery(field string) *elastic.ExistsQuery {
	return elastic.NewExistsQuery(field)
}

type BoolQueryAction string

var (
	BoolQueryMust    BoolQueryAction = "must"
	BoolQueryFilter  BoolQueryAction = "filter"
	BoolQueryShould  BoolQueryAction = "should"
	BoolQueryMustNot BoolQueryAction = "must_not"
)

// 文档 - 搜索条件 - bool 组合查找
// must			条件必须要满足，并将对分数起作用
// filter		条件必须要满足，但又不同于 must 子句，在 filter context 中执行，这意味着忽略评分，并考虑使用缓存。效率会高于 must
// should		条件应该满足。可以通过 minimum_should_match 参数指定应该满足的条件个数。如果 bool 查询包含 should 子句，并且没有 must 和 filter 子句，则默认值为 1，否则默认值为 0
// must_not		条件必须不能满足。在 filter context 中执行，这意味着评分被忽略，并考虑使用缓存。因为评分被忽略，所以会返回所有 0 分的文档
func (cli *client) DocSearchBoolQuery(action BoolQueryAction, queries ...elastic.Query) *elastic.BoolQuery {
	switch action {
	case BoolQueryMust:
		return elastic.NewBoolQuery().Must(queries...)

	case BoolQueryFilter:
		return elastic.NewBoolQuery().Filter(queries...)

	case BoolQueryShould:
		return elastic.NewBoolQuery().Should(queries...)

	case BoolQueryMustNot:
		return elastic.NewBoolQuery().MustNot(queries...)
	}

	return elastic.NewBoolQuery()
}

// 文档 - 查询文档数量
func (cli *client) DocCount(index string, query elastic.Query) (int64, error) {
	return cli.cli.Count(index).Query(query).Do(cli.ctx)
}
