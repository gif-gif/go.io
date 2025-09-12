package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取所有服务器的网卡总出站带宽 单位 byte/s, 如果不指定时间窗口，则默认获取90s内的数据来计算
func (g *GoPrometheus) GetSysTotalBandwidthOut(ctx context.Context, query MetricQuery) (model.Vector, error) {
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

	queryStr := fmt.Sprintf(`sum(rate(%s{%s}[%s]))`, MetricNodeTrafficOut, strings.Join(filters, ","), timeRange)

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysTotalBandwidthOut(vector *model.Vector, result map[int64]*Bandwidth) map[int64]*Bandwidth {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = &Bandwidth{}
		}
		result[instanceId].Out = int64(sample.Value)
	}
	return result
}

func (g *GoPrometheus) SysTotalBandwidthOut(ctx context.Context, query MetricQuery) (map[int64]*Bandwidth, error) {
	vector, err := g.GetSysTotalBandwidthOut(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*Bandwidth)
	result = g.PreHandleSysTotalBandwidthOut(&vector, result)
	return result, nil
}
