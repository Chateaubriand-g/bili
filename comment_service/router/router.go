package router

import (
	"github.com/Chateaubriand-g/bili/comment_service/controller"
	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
)

func InitRouter(ctl *controller.CommentController, tracer *zipkin.Tracer) *gin.Engine {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/comment/get", ctl.GetCommentList)
		v1.GET("/comment/add", ctl.AddComment)
		v1.GET("/comment/like", ctl.ClickLike)
	}

	return r
}
