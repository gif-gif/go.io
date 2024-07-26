package goeso

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// 文档 - 查询文档 - 根据ID
func (cli *GoEs) DocExists(index, id string) (bool, error) {
	return cli.es.Exists(index, id).IsSuccess(nil)
}

// 文档 - 搜索条件 - 判断某个字段是否存在
func (cli *GoEs) DocSearchExistsQuery(field string) types.Query {
	return types.Query{
		Exists: &types.ExistsQuery{Field: field},
	}
}

// 文档 - 查询文档 - 根据ID
func (cli *GoEs) DocGet(index, id string, m interface{}) (bool, error) {
	res, err := cli.es.Get(index, id).Refresh(true).Do(cli.ctx)
	if err != nil {
		return false, err
	}

	if m != nil {
		if len(res.Source_) > 0 {
			err = json.Unmarshal(res.Source_, m)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

// 文档 - 查询文档 - 根据ID集合
func (cli *GoEs) DocMultiGet(index string, ids []string, p *Pagination) (*search.Response, error) {
	return cli.DocSearch(index, types.Query{
		Ids: &types.IdsQuery{
			Values: ids,
		},
	}, p)
}

// query := `{ "query": { "match_all": {} } }`
func (cli *GoEs) DocSearchAll(index string, p *Pagination) (*search.Response, error) {
	return cli.DocSearch(index, types.Query{MatchAll: &types.MatchAllQuery{}}, p)
}

// 文档 - 搜索文档 - 根据条件
func (cli *GoEs) DocSearch(index string, query types.Query, p *Pagination) (*search.Response, error) {
	if p == nil {
		p = &Pagination{
			Sort:   "",
			Limit:  1000,
			Offset: 0,
		}
	}

	s := cli.es.Search()
	if index != "" {
		s.Index(index)

	}
	s.Request(&search.Request{
		Query: &query,
	}).Size(p.Limit).From(p.Offset)

	if p.Sort != "" {
		s.Sort(p.Sort)
	}
	res, err := s.Do(cli.ctx)

	if err != nil {
		return nil, err
	}

	// 遍历所有结果
	for _, hit := range res.Hits.Hits {
		fmt.Printf("%s\n", hit.Source_)
	}

	return res, nil
}

// 文档 - 搜索条件 - 精确匹配单个字段
func (cli *GoEs) DocSearchTermQuery(field string, value any) types.Query {
	return types.Query{
		Term: map[string]types.TermQuery{
			field: {Value: value},
		},
	}
}

// 文档 - 搜索条件 - 精确匹配多个字段
func (cli *GoEs) DocSearchTermsQuery(field string, values []interface{}) types.Query {
	t := types.NewTermsQuery()
	t.TermsQuery = make(map[string]types.TermsQueryField)
	t.TermsQuery[field] = values
	return types.Query{
		Terms: t,
	}
}

// 文档 - 搜索条件 - 匹配查找，单字段搜索（匹配分词结果）
func (cli *GoEs) DocSearchMatchQuery(field string, value string) types.Query {
	return types.Query{
		Match: map[string]types.MatchQuery{field: types.MatchQuery{Query: value}},
	}
}

// 文档 - 搜索条件 - 脚本查询
func (cli *GoEs) DocSearchScriptQuery(script types.Script) types.Query {
	return types.Query{
		Script: &types.ScriptQuery{
			Script: script,
		},
	}
}
func (cli *GoEs) DocSearchBoolQuery(terms map[string]interface{}) types.Query {
	ts := map[string]types.TermQuery{}
	for k, v := range terms {
		ts[k] = types.TermQuery{Value: v}
	}
	return types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				{Term: ts},
			},
		},
	}
}

// 文档 - 查询文档数量
//func (cli *GoEs) DocCount(index string, query elastic.Query) (int64, error) {
//	return cli.cli.Count(index).Query(query).Do(cli.ctx)
//}

// 文档 - 搜索条件 - 范围查找
//func (cli *GoEs) DocSearchRangeQuery(field string, from, to interface{}) *elastic.RangeQuery {
//	//return elastic.NewRangeQuery(field).Gte(from).Lte(to)
//
//	return types.Query{
//		Ra: elastic.NewRangeQuery(field).From(from).To(to),
//	}
//
//	elastic.NewRangeQuery(field).From(from).To(to)
//}

//
//type BoolQueryAction string
//
//var (
//	BoolQueryMust    BoolQueryAction = "must"
//	BoolQueryFilter  BoolQueryAction = "filter"
//	BoolQueryShould  BoolQueryAction = "should"
//	BoolQueryMustNot BoolQueryAction = "must_not"
//)
//
//// 文档 - 搜索条件 - bool 组合查找
//// must			条件必须要满足，并将对分数起作用
//// filter		条件必须要满足，但又不同于 must 子句，在 filter context 中执行，这意味着忽略评分，并考虑使用缓存。效率会高于 must
//// should		条件应该满足。可以通过 minimum_should_match 参数指定应该满足的条件个数。如果 bool 查询包含 should 子句，并且没有 must 和 filter 子句，则默认值为 1，否则默认值为 0
//// must_not		条件必须不能满足。在 filter context 中执行，这意味着评分被忽略，并考虑使用缓存。因为评分被忽略，所以会返回所有 0 分的文档
//func (cli *client) DocSearchBoolQuery(action BoolQueryAction, queries ...elastic.Query) *elastic.BoolQuery {
//	switch action {
//	case BoolQueryMust:
//		return elastic.NewBoolQuery().Must(queries...)
//
//	case BoolQueryFilter:
//		return elastic.NewBoolQuery().Filter(queries...)
//
//	case BoolQueryShould:
//		return elastic.NewBoolQuery().Should(queries...)
//
//	case BoolQueryMustNot:
//		return elastic.NewBoolQuery().MustNot(queries...)
//	}
//
//	return elastic.NewBoolQuery()
//}
//
//// 文档 - 查询文档数量
//func (cli *client) DocCount(index string, query elastic.Query) (int64, error) {
//	return cli.cli.Count(index).Query(query).Do(cli.ctx)
//}
