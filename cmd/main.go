package main

import (
	"fmt"
	"hezzl/config"
	"hezzl/internal/app"
	"hezzl/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.MustLoad()
	conf := config.GetConfig()

	logger.InitLog(logger.LogConfig{
		Mode:     conf.Env,
		LogPath:  conf.LogOutput,
		LogLevel: conf.LogLevel,
	})

	log := logger.GetLogger()
	app := app.New()
	appErr := make(chan error, 1)

	go func() {
		if err := app.Start(); err != nil {
			appErr <- fmt.Errorf("application failed to start: %w", err)
		}
	}()

	stop := make(chan os.Signal, 3)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case err := <-appErr:
		log.Error("application error", "error", err)
		if stopErr := app.Stop(); stopErr != nil {
			log.Error("additional error during shutdown", "error", stopErr)
		}
		os.Exit(1)
	case sig := <-stop:
		log.Info("received signal, shutting down", "signal", sig)
		if err := app.Stop(); err != nil {
			log.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
		log.Info("shutdown completed successfully")
	}
}
