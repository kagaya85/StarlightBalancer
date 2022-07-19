// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"starlight/services/timeline/internal/biz"
	"starlight/services/timeline/internal/conf"
	"starlight/services/timeline/internal/data"
	"starlight/services/timeline/internal/server"
	"starlight/services/timeline/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	httpServer := server.NewMetricServer(confServer, logger)
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	timelineRepo := data.NewTimelineRepo(dataData, logger)
	timelineUsecase := biz.NewTimelineUsecase(timelineRepo, logger)
	timelineService := service.NewTimelineService(timelineUsecase)
	grpcServer := server.NewGRPCServer(confServer, timelineService, logger)
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
