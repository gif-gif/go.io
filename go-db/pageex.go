package godb

//简单的分布组件
import "github.com/gogf/gf/util/gconv"

type OrderItem struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc"`
}

type Page struct {
	PageNo      int64        `json:"page_no,optional"`
	PageSize    int64        `json:"page_size,optional"`
	StartTime   int64        `json:"start_time,optional"`
	EndTime     int64        `json:"end_time,optional"`
	SortBy      []*OrderItem `json:"sort_by,optional"`
	GroupBy     []string     `json:"group_by,optional"`
	IgnoreTotal bool         `json:"ignore_total,optional"`
	IgnoreList  bool         `json:"need_list,optional"`
	OnlyTotal   bool         `json:"only_total,optional"`
	Ids         []int64      `json:"ids,optional"`
	States      []int64      `json:"states,optional"`
	Statuses    []int64      `json:"statutes,optional"`
}

func (p *Page) OrderBy() string {
	size := len(p.SortBy)
	if size == 0 {
		return "order by id desc"
	}

	order := "order by "
	for i, v := range p.SortBy {
		order = order + v.Column + " "
		if !v.Asc {
			order = order + " desc "
		}
		if size-1 == i {

		} else {
			order = order + ","
		}
	}
	return order
}

func (p *Page) OrderByExt() string {
	size := len(p.SortBy)
	if size == 0 {
		return ""
	}

	order := ""
	for i, v := range p.SortBy {
		order = order + v.Column + " "
		if !v.Asc {
			order = order + " desc "
		}
		if size-1 == i {

		} else {
			order = order + ","
		}
	}
	return order
}

func (p *Page) PageLimit() string {
	if p.PageNo == 0 {
		p.PageNo = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 20
	}
	return "limit " + gconv.String((p.PageNo-1)*p.PageSize) + "," + gconv.String(p.PageSize)
}

func (p *Page) GroupByStr() string {
	size := len(p.GroupBy)
	if size == 0 {
		return ""
	}

	order := "group by "
	for i, v := range p.GroupBy {
		order = order + v + " "
		if size-1 == i {

		} else {
			order = order + ","
		}
	}
	return order
}
