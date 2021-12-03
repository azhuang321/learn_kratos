// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"chat/app/group/service/internal/biz"
	"chat/app/group/service/internal/conf"
	"chat/app/group/service/internal/data"
	"chat/app/group/service/internal/server"
	"chat/app/group/service/internal/service"
	"chat/pkg/registry"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger, *zap.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, biz.ProviderSet, data.ProviderSet, registry.ProviderSet, newApp))
}
