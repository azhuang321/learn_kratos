package registry

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func NewConsulRegistry(logger *zap.Logger) registry.Registrar {
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		logger.Fatal(err.Error())
	}
	return consul.New(consulClient)
}
