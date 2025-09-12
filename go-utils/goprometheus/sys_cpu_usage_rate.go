package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
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

func (g *GoPrometheus) PreHandleSysCpuUsage(vector *model.Vector, result map[int64]*SysUsage) map[int64]*SysUsage {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = &SysUsage{}
		}
		result[instanceId].CpuUsage = float64(sample.Value)
	}

	return result
}

func (g *GoPrometheus) SysCpuUsageRate(ctx context.Context, query MetricQuery) (model.Vector, error) {
	vector, err := g.GetSysCpuUsageRate(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*SysUsage)
	result = g.PreHandleSysCpuUsage(&vector, result)
	return vector, nil
}
