package grpc

import (
	"hosting-resources-service/cmd/server/grpc/handlers/poolgrp"
	"hosting-resources-service/internal/pool"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	PoolBus *pool.Business
}

type App struct {
	server *grpc.Server
}

func New(cfg Config) *App {
	gs := grpc.NewServer()

	poolgrp.Register(gs, cfg.PoolBus)

	return &App{
		server: gs,
	}
}

func (a *App) Serve(lis net.Listener) error {
	return a.server.Serve(lis)
}

func (a *App) Stop() {
	a.server.GracefulStop()
}
