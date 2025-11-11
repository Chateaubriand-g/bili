package main

import (
	"log"

	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/notification_service/config"
	"github.com/Chateaubriand-g/bili/notification_service/controller"
	"github.com/Chateaubriand-g/bili/notification_service/dao"
	"github.com/Chateaubriand-g/bili/notification_service/util"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	tracer, reporter, err := middleware.InitZipkin(cfg)
	if err != nil {
		log.Fatalf("init zipkin error: %v", err)
	}
	defer middleware.CloseZipkin(reporter)

	deregister, err := middleware.RegisterServiceToConsul(cfg)
	if err != nil {
		log.Fatalf("register server to consul error: %v", err)
	}
	defer deregister()

	rds, err := middleware.InitRedis(cfg)
	if err != nil {
		log.Fatalf("init redis error: %v", err)
	}

	db, err := util.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("init databse err: %v", err)
	}

	userDAO := dao.NewNotifyDAO(db, rds)
	userController := controller.NewNotifyController(&userDAO)

	middleware.RegisterConsumer(cfg, db, rds)
	r := gin.Default()

	if tracer != nil {
		r.Use(middleware.ZipkinMiddleware(tracer))
	}

	api := r.Group("/api")
	{
		api.GET("/msg-unread/all", userController.GetUnreadByType)
	}

	r.Run(":8083")
}
