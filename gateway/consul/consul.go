package consul

import (
	"bili/gateway/config"

	"github.com/hashicorp/consul/api"
)

func NewConsul(cfg *config.Config) (*api.Client, error) {
	apicfg := api.DefaultConfig()
	apicfg.Address = cfg.Consul.Addr
	if cfg.Consul.Token != "" {
		apicfg.Token = cfg.Consul.Token
	}
	if cfg.Consul.Scheme != "" {
		apicfg.Scheme = cfg.Consul.Scheme
	}
	if cfg.Consul.Timeout > 0 {
		apicfg.HttpClient.Timeout = cfg.Consul.Timeout
	}

	return api.NewClient(apicfg)
}
