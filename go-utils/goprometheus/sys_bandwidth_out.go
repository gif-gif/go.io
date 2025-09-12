package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取各个服务器的总出站带宽 单位 byte/s, 如果不指定时间窗口，则默认获取90s内的数据来计算
func (g *GoPrometheus) GetSysBandwidthOut(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`device!~"tap.*|veth.*|br.*|docker.*|virbr*|lo*|cni.*|ifb.*"`,
	}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	timeRange := query.TimeRange
	if timeRange == "" {
		timeRange = "90s"
	}

	queryStr := fmt.Sprintf(`max by (%s) (rate(%s{%s}[%s]))`, MetricLabelInstanceId, MetricNodeTrafficOut, strings.Join(filters, ","), timeRange)

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysBandwidthOut(vector *model.Vector, result map[int64]*Bandwidth) map[int64]*Bandwidth {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = &Bandwidth{}
		}
		result[instanceId].Out = int64(sample.Value)
	}
	return result
}

func (g *GoPrometheus) SysBandwidthOut(ctx context.Context, query MetricQuery) (map[int64]*Bandwidth, error) {
	vector, err := g.GetSysBandwidthOut(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*Bandwidth)
	result = g.PreHandleSysBandwidthOut(&vector, result)
	return result, nil
}
