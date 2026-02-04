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
	oldLoginUseCase controller.OldLogin,
	oldRegisterUseCase controller.OldRegister,
	oldRefreshUseCase controller.OldRefreshTokens,
	oldLogoutUseCase controller.OldLogout,
	profileUseCase controller.UserProfile,
	editProfileUseCase controller.EditProfile,
	authorizeUseCase controller.Authorize,
	loginUseCase controller.Login,
	tokenUseCase controller.Token,
) *App {
	gRPCServer := grpcserver.New(grpcserver.WithPort(port))

	grpc.NewRouter(
		gRPCServer.App,
		oldLoginUseCase,
		oldRegisterUseCase,
		oldRefreshUseCase,
		oldLogoutUseCase,
		profileUseCase,
		editProfileUseCase,
		authorizeUseCase,
		loginUseCase,
		tokenUseCase,
	)

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
