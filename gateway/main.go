package main

import (
	"log"
	"fmt"

	"github.com/Chateaubriand-g/bili/gateway/config"
	"github.com/Chateaubriand-g/bili/gateway/consul"
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

	r := router.InitRouter(cli,cfg)
	addr := fmt.Sprintf(":%s",cfg.Gateway.Addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("gatway starting failed: %v", err)
	}
}
