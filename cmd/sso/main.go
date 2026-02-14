package main

import (
	"github.com/p1xray/pxr-sso/internal/app"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/pkg/logger"
	"log/slog"
)

func main() {
	cfgLoader := config.NewLoader()
	cfg := cfgLoader.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	application := app.New(log, cfg)

	go func() {
		application.Start()
	}()

	application.GracefulStop()
	log.Info("application stopped")
}
