package router

import (
	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/interaction_service/controller"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
)

func InitRouter(ctl *controller.InteractionController, tracer *zipkin.Tracer) {
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	v1 := r.Group("/api/v1")
	{
		videos := v1.Group("/videos")
		{
			videos.POST("/:video_id/likes", ctl.Like)
			videos.DELETE("/:video_id/likes", ctl.Like)
		}

		folders := v1.Group("/folders")
		{
			folders.POST("", ctl.CreateFolder)
			folders.DELETE("/:folder_id", ctl.DeleteFromFolder)

			folderVideos := folders.Group("/:folder_id/videos")
			{
				folderVideos.POST("/:video_id", ctl.AddToFolder)
				folderVideos.DELETE("/:video_id", ctl.DeleteFromFolder)
			}
		}
	}
}
