package main

import (
	"mgarnier11.fr/go/libs/logger"

	"mgarnier11.fr/go/go-autosaver/server"
)

func main() {
	logger.InitAppLogger("")

	api := server.NewServer(3000)

	api.Start()

}
