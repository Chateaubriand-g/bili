package router

import (
	"github.com/Chateaubriand-g/bili/auth_service/controller"
	"github.com/Chateaubriand-g/bili/common/middleware"
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

	v1 := r.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", ctl.Register)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", ctl.Login)
			auth.DELETE("/logout", ctl.Logout)
		}
	}

	r.GET("/swagger/*proxy", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
