package main

import (
	"mgarnier11.fr/go/libs/logger"

	"mgarnier11.fr/go/go-autosaver/config"
	"mgarnier11.fr/go/go-autosaver/server"
)

func main() {
	logger.InitAppLogger("GO-AUTOSAVER")

	api := server.NewServer(config.Config.ServerPort)
	api.Start()
}
