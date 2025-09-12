package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器cpu核心数（不是物理核心数，而是逻辑核心数。开启超线程后，一个物理核心可以对应多个逻辑核心，当前的指标中无法得知是否开启超线程，也无法区分物理核心和逻辑核心）
func (g *GoPrometheus) GetSysCpuCoreCount(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{}
	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	queryStr := fmt.Sprintf(`count by (%s) (count by (%s,%s) (%s{%s}))`, MetricLabelInstanceId, MetricLabelInstanceId, MetricLabelCpu, MetricNodeCpuSecondsTotal, strings.Join(filters, ","))

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysCpuCoreCount(vector *model.Vector, result map[int64]int64) map[int64]int64 {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		result[instanceId] = int64(sample.Value)
	}

	return result
}

func (g *GoPrometheus) SysCpuCoreCount(ctx context.Context, query MetricQuery) (map[int64]int64, error) {
	vector, err := g.GetSysCpuCoreCount(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int64)
	result = g.PreHandleSysCpuCoreCount(&vector, result)
	return result, nil
}
