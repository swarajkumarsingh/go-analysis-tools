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

var (
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "app_total_requests",
			Help: "Total number of requests to my app",
		},
	)
)

var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "cpu_temperature_celsius_a",
	Help: "Current temperature of the CPU.",
})

func CustomMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cpuTemp.Set(float64(100))
		totalRequests.Inc()
		fmt.Println(totalRequests.Desc())
		fmt.Println(totalRequests)
		c.Next()
	}
}

func init() {
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(totalRequests)
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
	r.Use(CustomMetricsMiddleware())
	r.Use(gintrace.Middleware("my-web-app"))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"error":   false,
			"message": "Health OK",
		})
	})

	r.GET("/big-task", func(ctx *gin.Context) {

		num := 0
		for i := 0; i < 10000000000; i++ {
			num = i
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"error":   false,
			"message": "OK",
			"number":  num,
		})
	})

	r.GET("/metrics", prometheusHandler())

	r.Run()
}
