package app

import (
	grpcapp "github.com/p1xray/pxr-sso/internal/app/grpc"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/redis"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/repository"
	"github.com/p1xray/pxr-sso/internal/oauth/infrastructure/storage/postgresql"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/authorize"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/consent"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/login"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/register"
	"github.com/p1xray/pxr-sso/internal/oauth/usecase/token"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// App is an application.
type App struct {
	log     *slog.Logger
	grpcApp *grpcapp.App
}

// New creates a new application.
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// Storages.
	// TODO: get postgresql connection URL from config
	storage, err := postgresql.New("postgresql://postgres:admin@127.0.0.1:5432/sso?sslmode=disable")
	if err != nil {
		panic(err)
	}

	// TODO: get redis connection URL from config
	redisStorage, err := redis.New("redis://test_redis_user:test_redis_user_pass@localhost:6380/0")
	if err != nil {
		panic(err)
	}

	// Repositories.
	oauthRepository := repository.NewRepository(storage)

	// Use-cases.
	authorizeUseCase := authorize.New(log, oauthRepository, redisStorage)
	loginUseCase := login.New(log, oauthRepository, redisStorage)
	registerUseCase := register.New(log, oauthRepository, redisStorage)
	consentUseCase := consent.New(log, oauthRepository, redisStorage)
	tokenUseCase := token.New(log, oauthRepository, redisStorage)

	grpcApp := grpcapp.New(
		log,
		cfg.GRPC.Port,
		authorizeUseCase,
		loginUseCase,
		registerUseCase,
		consentUseCase,
		tokenUseCase,
	)

	return &App{
		log:     log,
		grpcApp: grpcApp,
	}
}

// Start - starts the application.
func (a *App) Start() {
	const op = "app.Start"

	log := a.log.With(slog.String("op", op))
	log.Info("starting application")

	a.grpcApp.Start()
}

// GracefulStop - gracefully stops the application.
func (a *App) GracefulStop() {
	const op = "app.GracefulStop"

	log := a.log.With(slog.String("op", op))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-stop:
		log.Info("signal received from OS", slog.String("signal:", s.String()))
	case err := <-a.grpcApp.Notify():
		log.Error("received an error from the gRPC server:", sl.Err(err))
	}

	log.Info("stopping application")

	a.grpcApp.Stop()
}
