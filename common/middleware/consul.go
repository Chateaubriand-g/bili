package middleware

import (
	"fmt"
	"log"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/hashicorp/consul/api"
)

func RegisterServiceToConsul(cfg *config.Config) (func(), error) {
	consulcfg := api.DefaultConfig()

	cfgAddr := fmt.Sprintf("%s:%s", cfg.Consul.Addr, cfg.Consul.Port)

	consulcfg.Address = cfgAddr
	client, err := api.NewClient(consulcfg)
	if err != nil {
		return nil, fmt.Errorf("initial consul client failed: %w", err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      cfg.Server.ID,
		Name:    cfg.Server.Name,
		Address: cfg.Server.Addr,
		Port:    cfg.Server.Port,
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return nil, fmt.Errorf("service register failed: %w", err)
	}

	deregister := func() {
		if err := client.Agent().ServiceDeregister(cfg.Server.ID); err != nil {
			log.Printf("service deregister failed: %v", err)
			return
		}
		log.Printf("service deregister successful: %v", err)
	}

	return deregister, nil
}
