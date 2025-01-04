package main

import (
	"mgarnier11/go/logger"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/database"
	"mgarnier11/mineager/server"
	"mgarnier11/mineager/server/models"
)

func initDatabase() {
	database.InitDB()

	models.InitMapTable()
}

func main() {
	logger.InitAppLogger("mineager")

	initDatabase()

	api := server.NewServer(config.Config.ServerPort)

	api.Start()
}
