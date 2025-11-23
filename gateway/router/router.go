package router

import (
	"github.com/Chateaubriand-g/bili/gateway/config"
	"github.com/Chateaubriand-g/bili/gateway/middleware"
	"github.com/Chateaubriand-g/bili/gateway/proxy"
	"github.com/openzipkin/zipkin-go"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func InitRouter(cli *api.Client, cfg *config.Config, tracer *zipkin.Tracer, rds *redis.Client) *gin.Engine {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	api := r.Group("/api")
	{
		refresh := api.Group("/refresh")
		refresh.Use(middleware.JWTRAuth(cfg, rds))
		{
			refresh.POST("token", proxy.ReverseProxy(cli, "auth-service", tracer))
		}

		auth := api.Group("/auth")
		{
			account := auth.Group("/account")
			{
				account.Any("/*proxy", proxy.ReverseProxy(cli, "auth-service", tracer))
			}

			//user.Use(middleware.JWTAuth(cfg))
			//user.Any("/*proxy",proxy.ReverseProxy(cli,"user_service"))
		}

		user := api.Group("/user")
		user.Use(middleware.JWTAuth(cfg, rds))
		{
			user.Any("/*proxy", proxy.ReverseProxy(cli, "user-service", tracer))
		}

		msg := api.Group("/msg")
		msg.Use(middleware.JWTAuth(cfg, rds))
		{
			msg.GET("/*proxy", proxy.ReverseProxy(cli, "notify-service", tracer))
		}

	}
	return r
}
