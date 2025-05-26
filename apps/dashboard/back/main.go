package main

import (
	"mgarnier11.fr/go/dashboard/config"
	"mgarnier11.fr/go/dashboard/server"
	"mgarnier11.fr/go/libs/logger"
)

func main() {
	logger.InitAppLogger("dashboard")

	api := server.NewServer(config.Config.ServerPort)

	api.Start()
}
