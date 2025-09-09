package main

import (
	"context"
	"github.com/p1xray/pxr-sso/internal/app"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/pkg/logger"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	application := app.New(log, cfg)

	go func() {
		application.Start(context.Background())
	}()

	application.GracefulStop()
	log.Info("application stopped")
}
