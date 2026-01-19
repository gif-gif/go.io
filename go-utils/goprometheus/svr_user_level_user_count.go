package goprometheus

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/util/gconv"
	"github.com/prometheus/common/model"
)

func (g *GoPrometheus) GetSvrUserLevelUserCount(ctx context.Context, metrics string, query MetricQuery) (model.Vector, error) {
	filters := []string{}
	filters = append(filters, g.Filters...)

	filters = ToGroupFilter(filters, query.Group)
	filters = ToInstanceIdsFilter(filters, query.InstanceIds)

	queryStr := fmt.Sprintf(`%s{%s}`, metrics, strings.Join(filters, ","))

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSvrUserLevelUserCount(vector *model.Vector, result map[int64][]*UserLevelUserCount) map[int64][]*UserLevelUserCount {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		userLevel := gconv.Int(string(sample.Metric[MetricLabelUserLevel]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = []*UserLevelUserCount{}
		}
		result[instanceId] = append(result[instanceId], &UserLevelUserCount{
			Level: int64(userLevel),
			Count: int64(sample.Value),
		})
	}
	return result
}

// 获取每个会员等级的用户数(包含转发)
func (g *GoPrometheus) SvrUserLevelUserCount(ctx context.Context, query MetricQuery) (map[int64][]*UserLevelUserCount, error) {
	vector, err := g.GetSvrUserLevelUserCount(ctx, MetricUserLevelUserCount, query)
	if err != nil {
		return nil, err
	}
	return g.PreHandleSvrUserLevelUserCount(&vector, map[int64][]*UserLevelUserCount{}), nil
}

// 获取每个会员等级的用户数(不包含转发)
func (g *GoPrometheus) SvrRealUserLevelUserCount(ctx context.Context, query MetricQuery) (map[int64][]*UserLevelUserCount, error) {
	vector, err := g.GetSvrUserLevelUserCount(ctx, MetricRealUserLevelUserCount, query)
	if err != nil {
		return nil, err
	}
	return g.PreHandleSvrUserLevelUserCount(&vector, map[int64][]*UserLevelUserCount{}), nil
}
