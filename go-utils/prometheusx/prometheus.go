package prometheusx

import (
	"fmt"
	"github.com/gif-gif/go.io/go-utils/prometheusx/metric"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	LevelWarning = "warn"
	LevelError   = "error"
	LevelPanic   = "panic"
)

var (
	metricServerErrorTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "server",
		Subsystem: "alert",
		Name:      "alert_count",
		Help:      "server error",
		Labels:    []string{"level", "module", "name"},
	})

	sharedCache = cache.New(time.Minute, 5*time.Minute)
)

func Init(config Config) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	StartAgent(addr, config.Path)
}

func Alert(level, module, name string) error {
	key := level + module + name
	if _, ok := sharedCache.Get(key); ok {
		return fmt.Errorf("alert too fast:%s", key)
	}

	if len(name) > 50 {
		return fmt.Errorf("name is too long,max is 50, %s", name)
	}

	metricServerErrorTotal.Inc(level, module, name)
	sharedCache.Set(key, struct{}{}, time.Second*2)
	return nil
}

func AlertErr(module, name string) {
	Alert(LevelError, module, name)
}

func AlertWarn(module, name string) {
	Alert(LevelWarning, module, name)
}

func AlertPanic(module, name string) {
	Alert(LevelPanic, module, name)
}
