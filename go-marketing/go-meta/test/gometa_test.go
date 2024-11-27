package main

import (
	"context"
	gocontext "github.com/gif-gif/go.io/go-context"
	golog "github.com/gif-gif/go.io/go-log"
	gometa "github.com/gif-gif/go.io/go-marketing/go-meta"
	"github.com/gif-gif/go.io/goio"
	"testing"
)

func TestAdmobApps(t *testing.T) {
	goio.Init(goio.DEVELOPMENT)
	err := gometa.Init(context.Background(), gometa.Config{
		ClientId:     "",
		ClientSecret: "",
		AccessToken:  "",
	})
	if err != nil {
		golog.Error(err)
		return
	}
	// 测试获取应用列表
	req := gometa.AudienceDataRequest{
		AggregationPeriod: gometa.AGGREGATION_PERIOD_HOUR,
		Since:             "",
		Until:             "",
		Breakdowns:        []string{gometa.BREAKDOWN_COUNTRY, gometa.BREAKDOWN_PLACEMENT},
		Metrics:           gometa.DefaultAudienceMetrics,
		Filter: []gometa.AudienceFilter{
			{
				Field:    gometa.BREAKDOWN_PLACEMENT,
				Operator: "in",
				Values: []string{
					"",
				},
			},
		},
		Limit: 10,
	}
	res, err := gometa.Default().GetAudienceReport(&req, "")
	golog.Info(res)
	<-gocontext.WithCancel().Done()
}
