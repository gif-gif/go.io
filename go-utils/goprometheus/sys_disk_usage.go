package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"

	"github.com/prometheus/common/model"
)

// 获取服务器磁盘使用率 单位 %
func (g *GoPrometheus) GetSysDiskUsage(ctx context.Context, query MetricQuery) (model.Vector, error) {
	filters := []string{
		`fstype=~"ext.?|xfs"`,
		`mountpoint="/"`,
	}

	filters = append(filters, g.Filters...)

	// filters = toGroupFilter(filters, query.Group) // node-exporter 没有 group 标签
	filters = toInstanceIdsFilter(filters, query.InstanceIds)

	totalQuery := fmt.Sprintf(`%s{%s}`, MetricNodeDiskSize, strings.Join(filters, ","))
	freeQuery := fmt.Sprintf(`%s{%s}`, MetricNodeDiskFree, strings.Join(filters, ","))
	availableQuery := fmt.Sprintf(`%s{%s}`, MetricNodeDiskAvailable, strings.Join(filters, ","))
	queryStr := fmt.Sprintf(`max by(%s) ((%s-%s)*100/(%s+(%s-%s)))`, MetricLabelInstanceId, totalQuery, freeQuery, availableQuery, totalQuery, freeQuery)

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSysDiskUsage(vector *model.Vector, result map[int64]*SysUsage) map[int64]*SysUsage {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = &SysUsage{}
		}
		result[instanceId].DiskUsage = float64(sample.Value)
	}

	return result
}

func (g *GoPrometheus) SysDiskUsage(ctx context.Context, query MetricQuery) (map[int64]*SysUsage, error) {
	vector, err := g.GetSysDiskUsage(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*SysUsage)
	result = g.PreHandleSysDiskUsage(&vector, result)
	return result, nil
}
