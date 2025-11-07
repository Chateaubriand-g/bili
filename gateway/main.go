package main

import (
	"fmt"
	"log"

	"github.com/Chateaubriand-g/bili/gateway/config"
	"github.com/Chateaubriand-g/bili/gateway/consul"
	"github.com/Chateaubriand-g/bili/gateway/middleware"
	"github.com/Chateaubriand-g/bili/gateway/router"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	cli, err := consul.NewConsul(cfg)
	if err != nil {
		log.Fatalf("consul init failed: %v", err)
	}

	tracer, reporter, err := middleware.InitZipkin(cfg)
	if err != nil {
		log.Fatalf("initZipkin failed: %v", err)
	}
	defer middleware.CloseZipkin(reporter)

	r := router.InitRouter(cli, cfg, tracer)
	addr := fmt.Sprintf(":%s", cfg.Gateway.Addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("gatway starting failed: %v", err)
	}
}
