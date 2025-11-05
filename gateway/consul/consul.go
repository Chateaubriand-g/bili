package consul

import (
	"fmt"
	"net/http"

	"github.com/Chateaubriand-g/bili/gateway/config"

	"github.com/hashicorp/consul/api"
)

func NewConsul(cfg *config.Config) (*api.Client, error) {
	apicfg := api.DefaultConfig()

	addr := fmt.Sprintf("%s:%s",cfg.Consul.Addr,cfg.Consul.Port)
	apicfg.Address = addr
	if cfg.Consul.Token != "" {
		apicfg.Token = cfg.Consul.Token
	}
	if cfg.Consul.Scheme != "" {
		apicfg.Scheme = cfg.Consul.Scheme
	}
	
	if apicfg.HttpClient == nil {
		apicfg.HttpClient = &http.Client{}
	}
	apicfg.HttpClient.Timeout = cfg.Consul.Timeout

	return api.NewClient(apicfg)
}
