package godb

import (
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type SelectController[T int64 | string] struct {
	Values  []T  `json:"values,optional"`
	Exclude bool `json:"exclude,optional"`
}

func (c *SelectController[T]) ClickHouseWhere(column string) (string, []T) {
	if len(c.Values) == 0 {
		return "", nil
	}

	var whereString string
	if c.Exclude {
		whereString = " not in ? "
	} else {
		whereString = " in ? "
	}

	return " " + column + " " + whereString, c.Values
}

func (c *SelectController[T]) MysqlWhere(column string) (string, []any) {
	if len(c.Values) == 0 {
		return "", nil
	}

	var whereString string
	if c.Exclude {
		whereString = column + " not in  "
	} else {
		whereString = column + " in  "
	}

	conditions, params := GenerateSliceIn[T](c.Values)
	return whereString + conditions, params
}

func WhereIntArray[T int | int64 | int32](items []int64) string {
	if len(items) == 0 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString(" (")
	for i := 0; i < len(items)-1; i++ {
		builder.WriteString(gconv.String(items[i]))
		builder.WriteString(",")
	}
	builder.WriteString(gconv.String(items[len(items)-1]))
	builder.WriteString(") ")
	return builder.String()
}

func WhereStringArray(items []string) string {
	if len(items) == 0 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString(" (")
	for i := 0; i < len(items)-1; i++ {
		builder.WriteString("'" + items[i] + "',")
	}
	builder.WriteString("'" + gconv.String(items[len(items)-1]))
	builder.WriteString("') ")
	return builder.String()
}

func GenerateSliceIn[T any](srcItems []T) (string, []any) {
	if len(srcItems) == 0 {
		return "", nil
	}

	targetItems := make([]any, 0, len(srcItems))
	builder := strings.Builder{}
	builder.WriteString(" ( ")
	for _, item := range srcItems {
		builder.WriteString("?,")
		targetItems = append(targetItems, item)
	}

	targetString := builder.String()
	targetString = targetString[:len(targetString)-1]

	return targetString + " ) ", targetItems
}

func GenerateSliceInEx[T any](fieldName string, srcItems []T) (string, []any) {
	if len(srcItems) == 0 {
		return "", nil
	}

	targetItems := make([]any, 0, len(srcItems))
	builder := strings.Builder{}
	builder.WriteString(" ( ")
	for _, item := range srcItems {
		builder.WriteString("?,")
		targetItems = append(targetItems, item)
	}

	targetString := builder.String()
	targetString = targetString[:len(targetString)-1]

	return targetString + " ) ", targetItems
}

type NumberRangeController[T int64 | int | float64] struct {
	Min *T `json:"min,optional"`
	Max *T `json:"max,optional"`
}

func (c *NumberRangeController[T]) Where(column string) (string, []any) {
	if c.Min == nil && c.Max == nil {
		return "", []any{}
	}

	var whereString string
	args := make([]any, 0, 2)
	if c.Min != nil {
		whereString = " " + column + "  >= ? "
		args = append(args, *c.Min)
	}
	if c.Max != nil {
		if len(whereString) > 0 {
			whereString += " and " + column + "  < ? "
		} else {
			whereString = " " + column + "  < ? "
		}
		args = append(args, *c.Max)
	}

	return whereString, args
}

// 返回没有order by 前缀
func (p *Page) ClickHouseOrderByExt() string {
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
