package goprometheus

import (
	"context"
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取各个服务器的网卡总入站流量 单位 byte, 如果不指定时间窗口，则默认获取10分钟内的数据来计算
func (g *GoPrometheus) GetSysTrafficIn(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`device!~"tap.*|veth.*|br.*|docker.*|virbr*|lo*|cni.*|ifb.*"`,
	}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	timeRange := query.TimeRange
	if timeRange == "" {
		timeRange = "10m"
	}

	queryStr := fmt.Sprintf(`sum by (%s) (increase(%s{%s}[%s]))`, MetricLabelInstanceId, MetricNodeTrafficIn, strings.Join(filters, ","), timeRange)

	return g.PrometheusQuery(ctx, queryStr)
}
