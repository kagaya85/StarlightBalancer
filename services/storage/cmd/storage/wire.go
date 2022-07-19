// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"starlight/services/storage/internal/biz"
	"starlight/services/storage/internal/conf"
	"starlight/services/storage/internal/data"
	"starlight/services/storage/internal/server"
	"starlight/services/storage/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
