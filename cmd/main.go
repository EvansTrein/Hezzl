package main

import (
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

	app := app.New()

	go func() {
		if err := app.Start(); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 3)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-stop

	if err := app.Stop(); err != nil {
		logger.Error("application failed to stop")
		panic(err)
	}
}
