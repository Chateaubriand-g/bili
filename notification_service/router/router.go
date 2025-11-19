package router

import (
	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/notification_service/controller"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
)

func InitRouter(ctl *controller.NotifyController, tracer *zipkin.Tracer) *gin.Engine {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/msg-unread/all", ctl.GetUnreadByType)
	}

	return r
}
