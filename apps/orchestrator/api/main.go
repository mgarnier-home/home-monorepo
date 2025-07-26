package main

import (
	"os"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator-api/config"
	"mgarnier11.fr/go/orchestrator-api/server"
)

func main() {
	logger.InitAppLogger("dashboard")

	if config.Env.SSHPrivateKey != "" {
		// Create /ssh directory if it doesn't exist
		err := os.MkdirAll("ssh", 0700)
		if err != nil {
			panic(err)
		}
		// Create the id_rsa file with the content of SSH_PRIVATE_KEY
		err = os.WriteFile("ssh/ssh_private_key", []byte(config.Env.SSHPrivateKey), 0600)
		if err != nil {
			panic(err)
		}
	}

	api := server.NewServer(config.Env.ServerPort)

	api.Start()
}
