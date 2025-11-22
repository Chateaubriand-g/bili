package main

import (
	"log"

	"github.com/!chateaubriand-g/bili/common/mq"
	"github.com/Chateaubriand-g/bili/analytics_service/internal"
	"github.com/Chateaubriand-g/bili/common/middleware"
)

func main() {
	cfg, err := internal.LoadConfig()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	db, err := internal.InitDatabase(cfg)
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

	producer, err := mq.NewProducer(cfg, "notifications")
	if err != nil {
		log.Fatalf("init producer error: %v", err)
	}

	commentDAO := internal.NewAnalyticsDAO(db, rds)
	commentController := internal.NewAnalyticsController(commentDAO, producer)

	r := internal.InitRouter(commentController, tracer)
	r.Run(":/8085")
}
