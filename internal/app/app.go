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
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const componentTag = "[pxr-sso-app]"

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
	a.log.Debug(componentTag + " starting application")

	a.grpcApp.Start()
}

// GracefulStop - gracefully stops the application.
func (a *App) GracefulStop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	s := <-stop

	a.log.Debug(componentTag + " signal received from OS: " + s.String())
	a.log.Debug(componentTag + " stopping application")

	a.grpcApp.Stop()
}
