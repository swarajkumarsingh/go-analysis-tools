package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-analysis-tools/conf"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	ddprofiler "gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
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
		log.Panicf("error while startirn datadog profiler: %s", err)
	}
	defer ddprofiler.Stop()

	r := gin.Default()

	r.Use(gintrace.Middleware("my-web-app"))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{
			"error":   false,
			"message": "HELLO SIMU DUDE",
		})
	})

	r.Run(conf.PORT)
}
