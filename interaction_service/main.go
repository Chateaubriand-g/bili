package main

import (
	"log"

	"github.com/Chateaubriand-g/bili/common/middleware"
	"github.com/Chateaubriand-g/bili/common/mq"
	"github.com/Chateaubriand-g/bili/interaction_service/config"
	"github.com/Chateaubriand-g/bili/interaction_service/controller"
	"github.com/Chateaubriand-g/bili/interaction_service/dao"
	"github.com/Chateaubriand-g/bili/interaction_service/router"
	"github.com/Chateaubriand-g/bili/interaction_service/util"
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

	rds, err := middleware.InitRedis(cfg)
	if err != nil {
		log.Fatalf("init redis err: %v", err)
	}

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

	interactDAO := dao.NewInteractionDAO(db, rds)
	producer, err := mq.NewProducer(cfg, "notify")

	interacrController := controller.NewInteractionController(interactDAO, producer)

	r := router.InitRouter(interacrController, tracer)

	r.Run(":/8084")
}
