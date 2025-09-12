package goprometheus

import (
	"context"
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器cpu使用率 单位 %, 如果不指定时间窗口，则默认获取5分钟内的数据来计算
func (g *GoPrometheus) GetSysCpuUsageRate(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`mode="idle"`,
	}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	timeRange := query.TimeRange
	if timeRange == "" {
		timeRange = "5m"
	}

	queryStr := fmt.Sprintf(`(1 - avg(rate(%s{%s}[%s])) by (%s)) * 100`, MetricNodeCpuSecondsTotal, strings.Join(filters, ","), timeRange, MetricLabelInstanceId)

	return g.PrometheusQuery(ctx, queryStr)
}
