package main

import (
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-api/config"
	"mgarnier11.fr/go/orchestrator-api/server"
)

func main() {
	logger.InitAppLogger("dashboard")

	api := server.NewServer(config.Env.ServerPort)

	api.Start()
}
