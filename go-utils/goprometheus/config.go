package goprometheus

type Config struct {
	Name          string   `json:"name,optional" yaml:"Name"`
	PrometheusUrl string   `json:"prometheusUrl,optional" yaml:"PrometheusUrl" default:"0.0.0.0:9090"`
	Filters       []string `json:"filters,optional" yaml:"Filters"`
}

func (c *Config) GetUrl() string {
	if c.PrometheusUrl == "" {
		c.PrometheusUrl = "0.0.0.0:9090"
	}
	return c.PrometheusUrl
}
