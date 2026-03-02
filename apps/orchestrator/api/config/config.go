package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	ComposeDir    string
	ServerPort    int
	SSHPrivateKey string
	BinariesPath  string
}

func getEnv() (env *EnvConfig) {
	utils.InitEnvFromFile()

	env = &EnvConfig{
		ComposeDir:    utils.GetEnv("ORCHESTRATOR_COMPOSE_DIRECTORY", "/workspaces/home-config/compose"),
		ServerPort:    utils.GetEnv("ORCHESTRATOR_SERVER_PORT", 3000),
		SSHPrivateKey: utils.GetEnv("ORCHESTRATOR_SSH_PRIVATE_KEY", ""),
		BinariesPath:  utils.GetEnv("ORCHESTRATOR_BINARIES_PATH", "/dist"),
	}

	println("Using compose directory:", env.ComposeDir)

	return env
}

var Env *EnvConfig = getEnv()
