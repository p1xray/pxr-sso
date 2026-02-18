package grpcapp

import (
	"github.com/p1xray/pxr-sso/internal/controller"
	"github.com/p1xray/pxr-sso/internal/controller/grpc"
	"github.com/p1xray/pxr-sso/pkg/grpcserver"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
)

const componentTag = "[pxr-sso-grpc-app]"

// App is an gRPC controller application.
type App struct {
	log        *slog.Logger
	port       string
	gRPCServer *grpcserver.Server
}

// New creates new gRPC controller application.
func New(
	log *slog.Logger,
	cfg Config,
	authorizeUseCase controller.Authorize,
	loginUseCase controller.Login,
	registerUseCase controller.Register,
	consentUseCase controller.Consent,
	tokenUseCase controller.Token,
) *App {
	gRPCServer := grpcserver.New(grpcserver.WithPort(cfg.Port))

	grpc.NewRouter(
		gRPCServer.App,
		authorizeUseCase,
		loginUseCase,
		registerUseCase,
		consentUseCase,
		tokenUseCase,
	)

	return &App{
		log:        log,
		port:       cfg.Port,
		gRPCServer: gRPCServer,
	}
}

// Start - starts the gRPC controller application.
func (a *App) Start() {
	a.log.Debug(componentTag + " running gRPC server on port " + a.port)

	a.gRPCServer.Start()
	a.handleError()
}

// Stop - stops the gRPC controller application.
func (a *App) Stop() {
	a.log.Info(componentTag + " stopping gRPC server")

	a.gRPCServer.Stop()
}

func (a *App) handleError() {
	go func() {
		select {
		case err := <-a.gRPCServer.Notify():
			a.log.Warn(componentTag+" received an error from the gRPC server:", sl.Err(err))
		default:
		}
	}()
}
