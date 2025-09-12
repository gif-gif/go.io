package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器在线状态
func (g *GoPrometheus) GetSysUp(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	queryStr := fmt.Sprintf(`%s{%s}`, MetricUp, strings.Join(filters, ","))

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysUp(vector *model.Vector, result map[int64]int64) map[int64]int64 {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		result[instanceId] = int64(sample.Value)
	}

	return result
}

func (g *GoPrometheus) SysUp(ctx context.Context, query MetricQuery) (map[int64]int64, error) {
	vector, err := g.GetSysUp(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int64)
	result = g.PreHandleSysUp(&vector, result)
	return result, nil
}
