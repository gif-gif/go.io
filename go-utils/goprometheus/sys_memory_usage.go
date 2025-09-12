package goprometheus

import (
	"context"
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器内存使用率 单位 %
func (g *GoPrometheus) GetSysMemoryUsage(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	totalQuery := fmt.Sprintf(`%s{%s}`, MetricNodeMemTotal, strings.Join(filters, ","))
	availableQuery := fmt.Sprintf(`%s{%s}`, MetricNodeMemAvailable, strings.Join(filters, ","))
	queryStr := fmt.Sprintf(`(1 - (%s / %s)) * 100`, availableQuery, totalQuery)

	return g.PrometheusQuery(ctx, queryStr)
}
