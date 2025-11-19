package router

import (
	"github.com/Chateaubriand-g/bili/auth_service/controller"
	"github.com/Chateaubriand-g/bili/auth_service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(ctl *controller.AuthController, tracer *zipkin.Tracer) *gin.Engine {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	api := r.Group("/api")
	{
		api.POST("/auth/account/register", ctl.Register)
		api.POST("/auth/account/login", ctl.Login)
		api.POST("/auth/account/logout", ctl.Logout)
	}

	r.GET("/swagger/*proxy", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
