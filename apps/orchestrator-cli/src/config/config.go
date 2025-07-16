package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	OliveTinConfigDir string
	ComposeDir        string
	EnvDir            string
	SshKeyPath        string
}

func getEnv() (env *EnvConfig) {
	env = &EnvConfig{
		OliveTinConfigDir: utils.GetEnv("OLIVETIN_CONFIG_DIR", "/workspaces/home-config/olivetin"),
		ComposeDir:        utils.GetEnv("COMPOSE_DIR", "/workspaces/home-config/compose"),
		EnvDir:            utils.GetEnv("ENV_DIR", "/workspaces/home-config/compose"),
		SshKeyPath:        utils.GetEnv("SSH_KEY_PATH", ""),
	}

	return env
}

var Env *EnvConfig = getEnv()
