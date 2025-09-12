package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取所有服务器的网卡总入站流量 单位 byte, 如果不指定时间窗口，则默认获取10分钟内的数据来计算
func (g *GoPrometheus) GetSysTotalTrafficIn(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`device!~"tap.*|veth.*|br.*|docker.*|virbr*|lo*|cni.*|ifb.*"`,
	}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = ToInstanceIdsFilter(filters, query.InstanceIds)

	timeRange := query.TimeRange
	if timeRange == "" {
		timeRange = "10m"
	}

	queryStr := fmt.Sprintf(`sum(increase(%s{%s}[%s]))`, MetricNodeTrafficIn, strings.Join(filters, ","), timeRange)

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysTotalTrafficIn(vector *model.Vector, result map[int64]*Traffic) map[int64]*Traffic {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = &Traffic{}
		}
		result[instanceId].In = int64(sample.Value)
	}
	return result
}

func (g *GoPrometheus) SysTotalTrafficIn(ctx context.Context, query MetricQuery) (map[int64]*Traffic, error) {
	vector, err := g.GetSysTotalTrafficIn(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*Traffic)
	result = g.PreHandleSysTotalTrafficIn(&vector, result)
	return result, nil
}
