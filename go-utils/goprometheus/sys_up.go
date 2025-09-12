package goprometheus

import (
	"context"
	"fmt"
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
