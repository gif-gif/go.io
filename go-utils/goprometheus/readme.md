### Prometheus 系统指标
- [x] 基于 https://github.com/prometheus/node_exporter

一、 Labels 格式 `job=%s-node-exporter`

二、 查询系统指标过滤表达式实例，默认以产品分类
```
fmt.Sprintf(`%s="%s-node-exporter"`, MetricLabelJob, query.ProductCode)
```

三、用法 
```
g := goprometheus.NewGoPrometheus()
```
