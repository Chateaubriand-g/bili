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
		comments := v1.Group("/comments")
		{
			comments.GET("", ctl.GetCommentList)
			comments.POST("", ctl.AddComment)
			comments.POST("/:comments-id/likes", ctl.ClickLike)
			comments.DELETE("/:comments-id/likes/:user-id", ctl.ClickLike)
		}
	}

	return r
}
