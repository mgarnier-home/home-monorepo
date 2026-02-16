package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	OrchestratorApiUrl string
	ComposeDir         string
	Local              bool
}

func getEnv() (env *EnvConfig) {
	utils.InitEnvFromFile()

	env = &EnvConfig{
		OrchestratorApiUrl: utils.GetEnv("API_ORCHESTRATOR_URL", "http://localhost:3000"),
		ComposeDir:         utils.GetEnv("COMPOSE_DIRECTORY", ""),
		Local:              utils.GetEnv("ORCHESTRATOR_LOCAL", "false") == "true",
	}

	return env
}

var Env *EnvConfig = getEnv()
