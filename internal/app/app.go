package app

import (
	"log/slog"
	grpcapp "pxr-sso/internal/app/grpc"
	"pxr-sso/internal/config"
	clientcrud "pxr-sso/internal/logic/crud/client"
	sessioncrud "pxr-sso/internal/logic/crud/session"
	usercrud "pxr-sso/internal/logic/crud/user"
	"pxr-sso/internal/logic/service/auth"
	"pxr-sso/internal/logic/service/profile"
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

	userCRUD := usercrud.New(storage, storage, storage)
	clientCRUD := clientcrud.New(storage)
	sessionCRUD := sessioncrud.New(storage, storage)

	authService := auth.New(
		log,
		cfg.Tokens.AccessTokenTTL,
		cfg.Tokens.RefreshTokenTTL,
		userCRUD,
		clientCRUD,
		sessionCRUD,
	)

	profileService := profile.New(log, userCRUD)

	grpcApp := grpcapp.New(log, cfg.GRPC.Port, authService, profileService)

	return &App{
		GRPCServer: grpcApp,
	}
}
