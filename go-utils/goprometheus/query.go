package goprometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/common/model"
	"github.com/zeromicro/go-zero/core/logx"
)

func (g *GoPrometheus) PrometheusQuery(ctx context.Context, query string) (model.Vector, error) {
	logx.WithContext(ctx).Debugf("[prometheusQuery] executing query: %s", query)

	result, warnings, err := g.Api.Query(ctx, query, time.Now())
	if err != nil {
		logx.WithContext(ctx).Errorf("[prometheusQuery] prometheus query failed: %v", err)
		return nil, err
	}

	if len(warnings) > 0 {
		logx.WithContext(ctx).Infof("[prometheusQuery] prometheus query warnings: %s, warnings: %v", query, warnings)
	}

	vector, ok := result.(model.Vector)
	if !ok {
		err := fmt.Errorf("prometheus query result type unexpected: %v", result)
		logx.WithContext(ctx).Errorf("[prometheusQuery] %v", err)
		return nil, err
	}

	return vector, nil
}
