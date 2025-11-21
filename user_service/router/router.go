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
		users := v1.Group("/users")
		{
			users.GET("/:user-id", ctl.GetPersonalInfo)
			users.PATCH("/:user-id/info", ctl.UpdateInfo)
			users.PUT("/:user-id/avatar", ctl.UpdateAvatar)
		}
	}

	return r
}
