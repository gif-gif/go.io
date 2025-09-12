package goprometheus

import (
	"fmt"
	"strconv"
	"strings"
)

func toGroupFilter(filters []string, group string) []string {
	if group == "" {
		return filters
	} else {
		return append(filters, fmt.Sprintf(`%s="%s"`, MetricLabelGroup, group))
	}
}

func toInstanceIdsFilter(filters []string, instanceIds []int64) []string {
	if len(instanceIds) == 0 {
		return filters
	}
	strIds := make([]string, len(instanceIds))
	for i, id := range instanceIds {
		strIds[i] = strconv.Itoa(int(id))
	}
	return append(filters, fmt.Sprintf(`%s=~"%s"`, MetricLabelInstanceId, strings.Join(strIds, "|")))
}
