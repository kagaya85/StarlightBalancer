// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"starlight/services/upload/internal/biz"
	"starlight/services/upload/internal/conf"
	"starlight/services/upload/internal/data"
	"starlight/services/upload/internal/server"
	"starlight/services/upload/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	uploaderRepo := data.NewUploaderRepo(dataData, logger)
	uploaderUsecase := biz.NewUploaderUsecase(uploaderRepo, logger)
	uploaderService := service.NewUploaderService(uploaderUsecase)
	grpcServer := server.NewGRPCServer(confServer, uploaderService, logger)
	app := newApp(logger, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
