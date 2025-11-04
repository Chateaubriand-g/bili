package router

import (
	"bili/gateway/proxy"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func InitRouter(cli *api.Client) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		noAtuh := api.Group("/")
		{
			noAtuh.Any("/auth/*proxy", proxy.ReverseProxy(cli, "auth-service"))
		}

		auth := api.Group("/")
		{
			auth.Any("/user/*proxy", proxy.ReverseProxy(cli, "user-service"))
		}
	}
	return r
}
