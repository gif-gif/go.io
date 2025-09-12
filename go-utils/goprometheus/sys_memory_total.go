package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器内存总大小 单位 byte
func (g *GoPrometheus) GetSysMemoryTotal(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{}
	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = ToInstanceIdsFilter(filters, query.InstanceIds)

	queryStr := fmt.Sprintf(`%s{%s}`, MetricNodeMemTotal, strings.Join(filters, ","))

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysMemoryTotal(vector *model.Vector, result map[int64]*SysUsage) map[int64]*SysUsage {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = &SysUsage{}
		}
		result[instanceId].MemoryTotal = int64(sample.Value)
	}

	return result
}

func (g *GoPrometheus) SysMemoryTotal(ctx context.Context, query MetricQuery) (map[int64]*SysUsage, error) {
	vector, err := g.GetSysMemoryTotal(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make(map[int64]*SysUsage)
	result = g.PreHandleSysMemoryTotal(&vector, result)

	return result, nil
}
