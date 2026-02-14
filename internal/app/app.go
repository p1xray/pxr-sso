package app

import (
	grpcapp "github.com/p1xray/pxr-sso/internal/app/grpc"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/builder"
	"github.com/p1xray/pxr-sso/internal/oauth/domain/generator"
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
	cfg Config,
) *App {
	// storages
	storage, err := postgresql.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	redisStorage, err := redis.New(cfg.Redis)
	if err != nil {
		panic(err)
	}

	// repositories
	oauthRepository := repository.NewRepository(storage)

	// builders
	uriBuilder := builder.NewURI(cfg.URIBuilder)

	// generators
	tokenGenerator := generator.NewToken(cfg.Token)

	// use cases
	authorizeUseCase := authorize.New(log, uriBuilder, oauthRepository, redisStorage)
	loginUseCase := login.New(log, uriBuilder, oauthRepository, redisStorage)
	registerUseCase := register.New(log, uriBuilder, oauthRepository, redisStorage)
	consentUseCase := consent.New(log, uriBuilder, oauthRepository, redisStorage)
	tokenUseCase := token.New(log, tokenGenerator, oauthRepository, redisStorage)

	grpcApp := grpcapp.New(
		log,
		cfg.GRPC,
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
