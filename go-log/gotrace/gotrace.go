package gotrace

// TraceName represents the tracing name.

// A Config is a opentelemetry config.
type Config struct {
	TraceName string  `json:"traceName"`
	Name      string  `json:",optional"`
	Endpoint  string  `json:",optional"`
	Sampler   float64 `json:",default=1.0"`
	Batcher   string  `json:",default=jaeger,options=jaeger|zipkin"`
}
