package app

import (
	"log/slog"
	grpcapp "pxr-sso/internal/app/grpc"
	"pxr-sso/internal/service/auth"
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
	authService := auth.New(log)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
