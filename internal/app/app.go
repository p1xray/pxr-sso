package app

import (
	"github.com/p1xray/pxr-sso/pkg/grpcserver"
	"github.com/p1xray/pxr-sso/pkg/logger/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const logTag = "[pxr-sso-app]"

// App is the application structure.
type App struct {
	di         *diContainer
	grpcServer grpcserver.Server
}

// New creates a new application and initializes all dependencies through the DI container.
func New() *App {
	a := &App{
		di: newDIContainer(),
	}

	a.initDeps()

	return a
}

// Start launches the application.
func (a *App) Start() {
	log := a.di.Logger()

	log.Info(logTag + " starting application")
	log.Debug(logTag+" application configuration", slog.Any("config", a.di.Config()))

	a.grpcServer.Start()
	a.handleError()
}

// GracefulStop gracefully stops the application.
func (a *App) GracefulStop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	s := <-stop

	log := a.di.Logger()
	log.Debug(logTag + " signal received from OS: " + s.String())
	log.Info(logTag + " stopping application...")

	a.grpcServer.Stop()

	log.Info(logTag + " application stopped")
}

func (a *App) handleError() {
	go func() {
		for {
			select {
			case err := <-a.grpcServer.Notify():
				if err != nil {
					a.di.Logger().Error(logTag+" received an error from the gRPC server:", sl.Err(err))
				}
			default:
			}
		}
	}()
}

// initDeps sequentially calls initialization functions.
// If you need to add a new step (migration, metrics, etc.),
// you must add the function to the inits slice.
func (a *App) initDeps() {
	inits := []func(){
		a.initGRPCServer,
	}

	for _, fn := range inits {
		fn()
	}
}

// initGRPCServer initializes the GRPC server.
func (a *App) initGRPCServer() {
	a.grpcServer = a.di.GRPCServer()
}
