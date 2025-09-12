package goprometheusx

import (
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

var (
	once sync.Once
)

// StartAgent starts a prometheus agent.
func StartAgent(addr, path string) {
	if len(addr) < 3 {
		//端口可能冲突
		//addr = "0.0.0.0:9101"
		return
	}

	if path == "" {
		path = "/metrics"
	}

	once.Do(func() {
		goutils.AsyncFunc(func() {
			http.Handle(path, promhttp.Handler())
			golog.InfoF("Starting prometheus agent at %s%s", addr, path)

			if err := http.ListenAndServe(addr, nil); err != nil {
				golog.Error(err)
			}
			//endless.NewServer(addr, promhttp.Handler()).ListenAndServe()
		})
	})
}

func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func PrometheusBind(engine *gin.Engine, path string) {
	if path == "" {
		path = "/metrics"
	}
	engine.GET(path, PromHandler(promhttp.Handler()))
}
