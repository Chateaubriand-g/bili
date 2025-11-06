package router

import (
	"github.com/Chateaubriand-g/bili/gateway/proxy"
	"github.com/Chateaubriand-g/bili/gateway/config"
	//"github.com/Chateaubriand-g/bili/gateway/middleware"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func InitRouter(cli *api.Client,cfg *config.Config) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			account := user.Group("/account")
			{
				account.Any("/*proxy",proxy.ReverseProxy(cli,"auth-service"))
			}

			//user.Use(middleware.JWTAuth(cfg))
			//user.Any("/*proxy",proxy.ReverseProxy(cli,"user_service"))
		}

	}
	return r
}
