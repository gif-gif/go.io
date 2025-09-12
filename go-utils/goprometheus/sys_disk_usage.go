package goprometheus

import (
	"context"
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器磁盘使用率 单位 %
func (g *GoPrometheus) GetSysDiskUsage(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`fstype=~"ext.?|xfs"`,
		`mountpoint="/"`,
	}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	totalQuery := fmt.Sprintf(`%s{%s}`, MetricNodeDiskSize, strings.Join(filters, ","))
	freeQuery := fmt.Sprintf(`%s{%s}`, MetricNodeDiskFree, strings.Join(filters, ","))
	availableQuery := fmt.Sprintf(`%s{%s}`, MetricNodeDiskAvailable, strings.Join(filters, ","))
	queryStr := fmt.Sprintf(`max by(%s) ((%s-%s)*100/(%s+(%s-%s)))`, MetricLabelInstanceId, totalQuery, freeQuery, availableQuery, totalQuery, freeQuery)

	return g.PrometheusQuery(ctx, queryStr)
}
