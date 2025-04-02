package main

import (
	"mgarnier11.fr/go/mineager/config"
	"mgarnier11.fr/go/mineager/server"
	"mgarnier11.fr/go/mineager/server/database"

	"mgarnier11.fr/go/libs/logger"
)

func initDatabase() {
	database.InitDB()
	database.InitMapTable()
}

func main() {
	logger.InitAppLogger("mineager")

	initDatabase()

	api := server.NewServer(config.Config.ServerPort)

	api.Start()
}
