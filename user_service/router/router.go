package router

import (
	"github.com/Chateaubriand-g/bili/user_service/controller"
	"github.com/Chateaubriand-g/bili/user_service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
)

func InitRouter(ctl *controller.UserController, tracer *zipkin.Tracer) *gin.Engine {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/user/info/get-one", ctl.GetPersonalInfo)
		v1.POST("/user/info/update", ctl.UpdateInfo)
		v1.POST("/user/avatar/update", ctl.UpdateAvatar)
	}

	return r
}
