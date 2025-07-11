package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	ComposeDir        string
}

func getEnv() (env *EnvConfig) {
	env = &EnvConfig{
		ComposeDir:        utils.GetEnv("COMPOSE_DIR", "/workspaces/home-config/compose"),	}

	return env
}

var Env *EnvConfig = getEnv()
