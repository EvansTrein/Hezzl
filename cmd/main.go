package main

import (
	"hezzl/config"
	"hezzl/pkg/logs"
)

func main() {
	config.MustLoad()
	logs.InitLog(logs.LogConfig{
		Mode:     config.GetConfig().Env,
		LogPath:  config.GetConfig().LogOutput,
		LogLevel: config.GetConfig().LogLevel,
	})

}
