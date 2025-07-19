package main

import (
	"os/exec"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-api/config"
	"mgarnier11.fr/go/orchestrator-api/server"
)

func main() {
	logger.InitAppLogger("dashboard")

	if config.Env.SshKeyPath != "" {
		err := exec.Command("chmod", "600", config.Env.SshKeyPath).Run()

		if err != nil {
			panic(err)
		}
	}

	api := server.NewServer(config.Env.ServerPort)

	api.Start()
}
