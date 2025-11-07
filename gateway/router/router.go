package router

import (
	"github.com/Chateaubriand-g/bili/gateway/config"
	"github.com/Chateaubriand-g/bili/gateway/middleware"
	"github.com/Chateaubriand-g/bili/gateway/proxy"
	"github.com/openzipkin/zipkin-go"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func InitRouter(cli *api.Client, cfg *config.Config, tracer *zipkin.Tracer) *gin.Engine {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			account := user.Group("/account")
			{
				account.Any("/*proxy", proxy.ReverseProxy(cli, "auth-service", tracer))
			}

			//user.Use(middleware.JWTAuth(cfg))
			//user.Any("/*proxy",proxy.ReverseProxy(cli,"user_service"))
		}

	}
	return r
}
