package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"

	"github.com/Chateaubriand-g/bili/common/middleware"
)

func InitRouter(ctl *AnalyticsController, tracer *zipkin.Tracer) *gin.Engine {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	v1 := r.Group("/api/v1")
	{
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/video/:id", ctl.GetVideoStats)
		}
	}

	return r
}
