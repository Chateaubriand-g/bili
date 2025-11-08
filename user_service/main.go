package main

import (
	"log"

	"github.com/Chateaubriand-g/bili/user_service/config"
	"github.com/Chateaubriand-g/bili/user_service/controller"
	"github.com/Chateaubriand-g/bili/user_service/dao"
	"github.com/Chateaubriand-g/bili/user_service/middleware"
	"github.com/Chateaubriand-g/bili/user_service/ossclient"
	"github.com/Chateaubriand-g/bili/user_service/util"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("init config failed: %v", err)
	}

	tracer, reporter, err := middleware.InitZipkin(cfg)
	if err != nil {
		log.Fatalf("init zipkin failed: %v", err)
	}
	defer middleware.CloseZipkin(reporter)

	deregister, err := middleware.RegisterServiceToConsul(cfg)
	if err != nil {
		log.Fatalf("register service failed: %v", err)
	}
	defer deregister()

	db, err := util.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("init databse failed: %v", err)
	}

	oss := ossclient.InitOssClient(cfg)

	userDAO := dao.NewUserDAO(db)
	userController := controller.NewUserDAO(userDAO, oss)

	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	api := r.Group("/api")
	{
		api.GET("/user/info/get-one", userController.GetPersonalInfo)
		api.POST("/user/info/update", userController.UpdateInfo)
		api.POST("/user/avatar/update", userController.UpdateAvatar)
	}

	r.Run(":8082")

}
