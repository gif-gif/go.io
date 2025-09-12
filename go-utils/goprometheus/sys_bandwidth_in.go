package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取各个服务器的总入站带宽 单位 byte/s, 如果不指定时间窗口，则默认获取90s内的数据来计算
func (g *GoPrometheus) GetSysBandwidthIn(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`device!~"tap.*|veth.*|br.*|docker.*|virbr*|lo*|cni.*|ifb.*"`,
	}
	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = ToInstanceIdsFilter(filters, query.InstanceIds)

	timeRange := query.TimeRange
	if timeRange == "" {
		timeRange = "90s"
	}

	queryStr := fmt.Sprintf(`max by (%s)(rate(%s{%s}[%s]))`, MetricLabelInstanceId, MetricNodeTrafficIn, strings.Join(filters, ","), timeRange)

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysBandwidthIn(vector *model.Vector, result map[int64]*Bandwidth) map[int64]*Bandwidth {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = &Bandwidth{}
		}
		result[instanceId].In = int64(sample.Value)
	}
	return result
}

func (g *GoPrometheus) SysBandwidthIn(ctx context.Context, query MetricQuery) (map[int64]*Bandwidth, error) {
	vector, err := g.GetSysBandwidthIn(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*Bandwidth)
	result = g.PreHandleSysBandwidthIn(&vector, result)
	return result, nil
}
