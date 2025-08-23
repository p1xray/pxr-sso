package grpcapp

import (
	"fmt"
	"github.com/p1xray/pxr-sso/internal/controller"
	authcontroller "github.com/p1xray/pxr-sso/internal/controller/grpc/auth"
	profilecontroller "github.com/p1xray/pxr-sso/internal/controller/grpc/profile"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

// App is an gRPC application.
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC controller application.
func New(
	log *slog.Logger,
	port int,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	refreshUseCase controller.RefreshTokens,
	logoutUseCase controller.Logout,
	profileUseCase controller.UserProfile,
) *App {
	gRPCServer := grpc.NewServer()

	authcontroller.Register(
		gRPCServer,
		loginUseCase,
		registerUseCase,
		refreshUseCase,
		logoutUseCase,
	)

	profilecontroller.Register(gRPCServer, profileUseCase)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC controller and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC controller.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf((":%d"), a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC controller is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC controller.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op))
	a.log.Info("stopping gRPC controller", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
