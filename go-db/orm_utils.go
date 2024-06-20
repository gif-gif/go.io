package godb

import (
	"strings"

	"github.com/gogf/gf/util/gconv"
)

// 元素都转换成字符串比较
func IsInArray[T any](arr []T, target T) bool {
	for _, num := range arr {
		tt := gconv.String(target)
		tt = strings.ReplaceAll(tt, "`", "")
		if gconv.String(num) == tt {
			return true
		}
	}
	return false
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
