package gateway

import (
	"bili/gateway/config"
	"bili/gateway/consul"
	"bili/gateway/router"
	"log"
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
	r := router.InitRouter(cli)
	addr := cfg.Gateway.Addr
	if err := r.Run(addr); err != nil {
		log.Fatalf("gatway starting failed: %v", err)
	}
}
