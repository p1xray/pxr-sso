package app

import (
	grpcapp "github.com/p1xray/pxr-sso/internal/app/grpc"
	kafkaapp "github.com/p1xray/pxr-sso/internal/app/kafka"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/internal/infrastructure/kafka/handlers"
	"github.com/p1xray/pxr-sso/internal/infrastructure/repository"
	"github.com/p1xray/pxr-sso/internal/infrastructure/storage/sqlite"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/login"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/logout"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/refresh"
	"github.com/p1xray/pxr-sso/internal/usecase/auth/register"
	"github.com/p1xray/pxr-sso/internal/usecase/profile/card"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// App is an application.
type App struct {
	log      *slog.Logger
	grpcApp  *grpcapp.App
	kafkaApp *kafkaapp.App
}

// New creates a new application.
func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// Storages.
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	numberOfTopics := 1 // TODO: get this from config.
	kafkaApp := kafkaapp.New(log, numberOfTopics)

	// Handlers.
	registerHandler := handlers.NewUserHasRegistered(log, kafkaApp.Input())

	// Repositories.
	authRepository := repository.NewAuthRepository(log, storage)
	profileRepository := repository.NewProfileRepository(log, storage)

	// Use-cases.
	loginUseCase := login.New(log, cfg.Tokens, authRepository)
	registerUseCase := register.New(log, cfg.Tokens, authRepository, registerHandler)
	refreshUseCase := refresh.New(log, cfg.Tokens, authRepository)
	logoutUseCase := logout.New(log, cfg.Tokens, authRepository)

	profileUseCase := card.New(log, profileRepository)

	grpcApp := grpcapp.New(
		log,
		cfg.GRPC.Port,
		loginUseCase,
		registerUseCase,
		refreshUseCase,
		logoutUseCase,
		profileUseCase,
	)

	return &App{
		log:      log,
		grpcApp:  grpcApp,
		kafkaApp: kafkaApp,
	}
}

// Start - starts the application.
func (a *App) Start() {
	const op = "app.Start"

	log := a.log.With(slog.String("op", op))
	log.Info("starting application")

	a.kafkaApp.Start()
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
	a.kafkaApp.Stop()
}
