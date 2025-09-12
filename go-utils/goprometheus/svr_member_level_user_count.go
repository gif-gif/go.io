package goprometheus

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/prometheus/common/model"
	"strings"
)

func (g *GoPrometheus) GetSvrMemberLevelUserCount(ctx context.Context, metrics string, query MetricQuery) (model.Vector, error) {
	filters := []string{}
	filters = append(filters, g.Filters...)

	filters = ToGroupFilter(filters, query.Group)
	filters = ToInstanceIdsFilter(filters, query.InstanceIds)

	queryStr := fmt.Sprintf(`%s{%s}`, metrics, strings.Join(filters, ","))

	return g.PrometheusQuery(ctx, queryStr)
}

func (g *GoPrometheus) PreHandleSvrMemberLevelUserCount(vector *model.Vector, result map[int64][]*MemberLevelUserCount) map[int64][]*MemberLevelUserCount {
	for _, sample := range *vector {
		instanceId := gconv.Int64(string(sample.Metric[MetricLabelInstanceId]))
		memberLevel := gconv.Int(string(sample.Metric[MetricLabelMemberLevel]))
		if _, ok := result[instanceId]; !ok {
			result[instanceId] = []*MemberLevelUserCount{}
		}
		result[instanceId] = append(result[instanceId], &MemberLevelUserCount{
			Level: int64(memberLevel),
			Count: int64(sample.Value),
		})
	}
	return result
}

// 获取每个会员等级的用户数(包含转发)
func (g *GoPrometheus) SvrMemberLevelUserCount(ctx context.Context, query MetricQuery) (map[int64][]*MemberLevelUserCount, error) {
	vector, err := g.GetSvrMemberLevelUserCount(ctx, MetricMemberLevelUserCount, query)
	if err != nil {
		return nil, err
	}
	return g.PreHandleSvrMemberLevelUserCount(&vector, map[int64][]*MemberLevelUserCount{}), nil
}

// 获取每个会员等级的用户数(不包含转发)
func (g *GoPrometheus) SvrRealMemberLevelUserCount(ctx context.Context, query MetricQuery) (map[int64][]*MemberLevelUserCount, error) {
	vector, err := g.GetSvrMemberLevelUserCount(ctx, MetricRealMemberLevelUserCount, query)
	if err != nil {
		return nil, err
	}
	return g.PreHandleSvrMemberLevelUserCount(&vector, map[int64][]*MemberLevelUserCount{}), nil
}
