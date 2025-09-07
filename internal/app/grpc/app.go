package grpcapp

import (
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/controller/grpc"
	"github.com/p1xray/pxr-sso/pkg/grpcserver"
	"log/slog"
)

// App is an gRPC controller application.
type App struct {
	log        *slog.Logger
	port       string
	gRPCServer *grpcserver.Server
}

// New creates new gRPC controller application.
func New(
	log *slog.Logger,
	port string,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	refreshUseCase controller.RefreshTokens,
	logoutUseCase controller.Logout,
	profileUseCase controller.UserProfile,
) *App {
	gRPCServer := grpcserver.New(grpcserver.WithPort(port))

	grpc.NewRouter(
		gRPCServer.App,
		loginUseCase,
		registerUseCase,
		refreshUseCase,
		logoutUseCase,
		profileUseCase)

	return &App{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

// Start - starts the gRPC controller application.
func (a *App) Start() {
	const op = "grpcapp.Start"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)
	log.Info("running gRPC server")

	a.gRPCServer.Start()
}

// Stop - stops the gRPC controller application.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)
	log.Info("stopping gRPC server")

	a.gRPCServer.Stop()
}

// Notify - notifies about gRPC controller application errors.
func (a *App) Notify() <-chan error {
	return a.gRPCServer.Notify()
}
