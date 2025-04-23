package app

import (
	"log/slog"
	grpcapp "pxr-sso/internal/app/grpc"
)

// App is an application.
type App struct {
	GRPCServer *grpcapp.App
}

// New creates a new application.
func New(
	log *slog.Logger,
	grpcPort int,
) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
