package app

import (
	grpcapp "github.com/p1xray/pxr-sso/internal/app/grpc"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/infrastructure/repository"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/sqlite"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/login"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/logout"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/refresh"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/register"
	"log/slog"
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

	authRepository := repository.NewAuthRepository(log, storage)

	loginUseCase := login.New(log, cfg.Tokens, authRepository)
	registerUseCase := register.New(log, cfg.Tokens, authRepository)
	refreshUseCase := refresh.New(log, cfg.Tokens, authRepository)
	logoutUseCase := logout.New(log, cfg.Tokens, authRepository)

	grpcApp := grpcapp.New(
		log,
		cfg.GRPC.Port,
		loginUseCase,
		registerUseCase,
		refreshUseCase,
		logoutUseCase,
	)

	return &App{
		GRPCServer: grpcApp,
	}
}
