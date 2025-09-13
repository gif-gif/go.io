### Prometheus 系统指标
- [x] 基于 https://github.com/prometheus/node_exporter

一、 Labels 格式 `job=%s-node-exporter`

二、 查询系统指标过滤表达式实例，默认以产品分类
```
fmt.Sprintf(`%s="%s-node-exporter"`, MetricLabelJob, query.ProductCode)
```

三、用法 
```
package main

import (
	"context"
	"fmt"
	"github.com/gif-gif/go.io/go-utils/goprometheus"
)

func main() {
	product := "fkey"
	err := goprometheus.Init(goprometheus.Config{
		Name:          product,
		PrometheusUrl: "http://127.0.0.1:9091",
		Filters: []string{
			//fmt.Sprintf(`%s="%s-node-exporter"`, goprometheus.MetricLabelJob, product),
			fmt.Sprintf(`%s="%s-node"`, goprometheus.MetricLabelJob, product),
		},
	})

	if err != nil {
		panic(err)
	}

	query := goprometheus.MetricQuery{
		ProductCode: product,
		Group:       "all",
		InstanceIds: []int64{2015},
	}

	memberCount, err := goprometheus.GetClient(product).SvrRealMemberLevelUserCount(context.Background(), query)
	if err != nil {
		panic(err)
	}
	fmt.Println(memberCount)

}

```
