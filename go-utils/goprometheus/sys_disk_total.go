package goprometheus

import (
	"context"
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器磁盘总大小 单位 byte
func (g *GoPrometheus) GetSysDiskTotal(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`fstype=~"ext.?|xfs"`,
		`mountpoint="/"`,
	}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	queryStr := fmt.Sprintf(`%s{%s}`, MetricNodeDiskSize, strings.Join(filters, ","))
	return g.PrometheusQuery(ctx, queryStr)
}
