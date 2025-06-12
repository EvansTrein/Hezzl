package main

import (
	"hezzl/config"
	"hezzl/pkg/logger"
)

func main() {
	config.MustLoad()
	logger.InitLog(logger.LogConfig{
		Mode:     config.GetConfig().Env,
		LogPath:  config.GetConfig().LogOutput,
		LogLevel: config.GetConfig().LogLevel,
	})

}
