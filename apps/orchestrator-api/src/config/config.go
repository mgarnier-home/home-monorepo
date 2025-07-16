package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	ComposeDir string
	ServerPort int
}

func getEnv() (env *EnvConfig) {
	utils.InitEnvFromFile()

	env = &EnvConfig{
		ComposeDir: utils.GetEnv("COMPOSE_DIRECTORY", "/workspaces/home-config/compose"),
		ServerPort: utils.GetEnv("SERVER_PORT", 3000),
	}

	println("Using compose directory:", env.ComposeDir)

	return env
}

var Env *EnvConfig = getEnv()
