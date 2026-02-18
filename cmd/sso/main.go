package main

import (
	"github.com/p1xray/pxr-sso/internal/app"
	"github.com/p1xray/pxr-sso/internal/config"
	"github.com/p1xray/pxr-sso/pkg/logger"
	"log/slog"
)

const componentTag = "[pxr-sso-main]"

func main() {
	cfgLoader := config.NewLoader()
	cfg := cfgLoader.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Debug(componentTag+" starting application", slog.Any("config", cfg))

	application := app.New(log, cfg)

	go application.Start()

	application.GracefulStop()
	log.Debug(componentTag + " application stopped")
}
