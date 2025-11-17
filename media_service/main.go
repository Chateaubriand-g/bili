package main

import (
	"context"
	"log"

	"github.com/Chateaubriand-g/bili/common/alioss"
	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/media_service/config"
	"github.com/Chateaubriand-g/bili/media_service/controller"
	"github.com/Chateaubriand-g/bili/media_service/router"
	"github.com/Chateaubriand-g/bili/media_service/util"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	db, err := util.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("init databse err: %v", err)
	}

	//rds, err := middleware.InitRedis(cfg)

	deregister, err := middleware.RegisterServiceToConsul(cfg)
	if err != nil {
		log.Fatalf("register server to consul error: %v", err)
	}
	defer deregister()

	tracer, reporter, err := middleware.InitZipkin(cfg)
	if err != nil {
		log.Fatalf("init zipkin error: %v", err)
	}
	defer middleware.CloseZipkin(reporter)

	oss := alioss.NewOssClient(context.TODO(), cfg)
	sts, _ := alioss.NewSTSClient(context.TODO())

	mediaControler := controller.NewMediaController(db, oss, sts)

	r := router.InitRouter(mediaControler, tracer)

	r.Run(":/8084")
}
