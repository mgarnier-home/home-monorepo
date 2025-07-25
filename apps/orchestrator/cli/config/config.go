package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	OrchestratorApiUrl string
}

func getEnv() (env *EnvConfig) {
	utils.InitEnvFromFile()

	env = &EnvConfig{
		OrchestratorApiUrl: utils.GetEnv("API_ORCHESTRATOR_URL", "http://localhost:3000"),
	}

	return env
}

var Env *EnvConfig = getEnv()
