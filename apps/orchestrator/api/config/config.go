package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	ComposeDir   string
	ServerPort   int
	SshKeyPath   string
	BinariesPath string
}

func getEnv() (env *EnvConfig) {
	utils.InitEnvFromFile()

	env = &EnvConfig{
		ComposeDir:   utils.GetEnv("COMPOSE_DIRECTORY", "/workspaces/home-config/compose"),
		ServerPort:   utils.GetEnv("SERVER_PORT", 3000),
		SshKeyPath:   utils.GetEnv("SSH_KEY_PATH", ""),
		BinariesPath: utils.GetEnv("BINARIES_PATH", "/dist"),
	}

	println("Using compose directory:", env.ComposeDir)

	return env
}

var Env *EnvConfig = getEnv()
