// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"chat/app/user/service/internal/biz"
	"chat/app/user/service/internal/conf"
	"chat/app/user/service/internal/data"
	"chat/app/user/service/internal/server"
	"chat/app/user/service/internal/service"
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
