package goeso

type Pagination struct {
	Sort   string // 排序字段 "example --> published_at:desc "
	Offset int
	Limit  int // default 1000
}
