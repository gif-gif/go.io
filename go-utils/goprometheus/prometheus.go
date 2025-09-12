package goprometheus

import (
	prometheusApi "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type GoPrometheus struct {
	Client  prometheusApi.Client
	Api     v1.API
	Filters []string // 服务器过滤条件
}

func New(config Config) (*GoPrometheus, error) {
	client, err := prometheusApi.NewClient(prometheusApi.Config{
		Address: config.GetAddress(),
	})
	if err != nil {
		return nil, err
	}

	v1API := v1.NewAPI(client)

	return &GoPrometheus{
		Client: client,
		Api:    v1API,
	}, nil
}

func (g *GoPrometheus) AddFilters(filters ...string) {
	g.Filters = append(g.Filters, filters...)
}

func (g *GoPrometheus) GetFilters() []string {
	return g.Filters
}

func (g *GoPrometheus) SetFilters(filters ...string) {
	if len(g.Filters) == 0 {
		g.Filters = []string{}
	} else {
		g.Filters = filters
	}
}
