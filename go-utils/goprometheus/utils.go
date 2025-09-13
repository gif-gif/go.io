package goprometheus

import (
	"fmt"
	"strconv"
	"strings"
)

func ToGroupFilter(filters []string, group string) []string {
	if group == "" {
		return filters
	} else {
		return ToFilter(filters, string(MetricLabelGroup), group)
	}
}

// `%s=~"%s"` 包含关系
func ToInstanceIdsFilter(filters []string, instanceIds []int64) []string {
	if len(instanceIds) == 0 {
		return filters
	}
	strIds := make([]string, len(instanceIds))
	for i, id := range instanceIds {
		strIds[i] = strconv.Itoa(int(id))
	}
	return append(filters, fmt.Sprintf(`%s=~"%s"`, MetricLabelInstanceId, strings.Join(strIds, "|")))
}

// 通用精确过滤器（`%s="%s"`）
func ToFilter(filters []string, metrics string, value string) []string {
	if value == "" {
		return filters
	} else {
		return append(filters, fmt.Sprintf(`%s="%s"`, metrics, value))
	}
}
