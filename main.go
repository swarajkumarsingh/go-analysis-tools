package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-analysis-tools/conf"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	ddprofiler "gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius",
	Help: "Current temperature of the CPU.",
})

func init() {
	prometheus.MustRegister(cpuTemp)
}

func main() {

	if conf.ENV == conf.ENV_PROD {
		gin.SetMode(gin.ReleaseMode)

		// starting datadog tracer
		ddtracer.Start(
			ddtracer.WithAgentAddr(fmt.Sprintf("%s:8126", conf.DDAgentHost)),
			ddtracer.WithDogstatsdAddress(fmt.Sprintf("%s:8125", conf.DDAgentHost)),
			ddtracer.WithEnv(conf.ClientENV),
			ddtracer.WithRuntimeMetrics(),
			ddtracer.WithService(conf.DDServiceName),
		)
		defer ddtracer.Stop()

		// starting datadog profiler
		err := ddprofiler.Start(
			ddprofiler.WithAgentAddr(fmt.Sprintf("%s:8126", conf.DDAgentHost)),
			ddprofiler.WithEnv(conf.ClientENV),
			ddprofiler.WithService(conf.DDServiceName),
			ddprofiler.WithProfileTypes(
				ddprofiler.CPUProfile,
				ddprofiler.HeapProfile,
				ddprofiler.BlockProfile,
				ddprofiler.GoroutineProfile,
				ddprofiler.MetricsProfile,
				ddprofiler.MutexProfile,
			),
		)
		if err != nil {
			log.Panicf("error while starting datadog profiler: %s", err)
		}
		defer ddprofiler.Stop()
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gintrace.Middleware("my-web-app"))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"error":   false,
			"message": "Health OK",
		})
	})

	r.GET("/metrics", prometheusHandler())

	r.Run(conf.PORT)
}
