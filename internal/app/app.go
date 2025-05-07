package app

import (
	"log/slog"
	grpcapp "pxr-sso/internal/app/grpc"
	"pxr-sso/internal/config"
	"pxr-sso/internal/logic/service/auth"
	"pxr-sso/internal/storage/sqlite"
)

// App is an application.
type App struct {
	GRPCServer *grpcapp.App
}

// New creates a new application.
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, cfg.Tokens.AccessTokenTTL, cfg.Tokens.RefreshTokenTTL)

	grpcApp := grpcapp.New(log, cfg.GRPC.Port, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
