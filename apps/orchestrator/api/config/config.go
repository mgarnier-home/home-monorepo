package config

import (
	"mgarnier11.fr/go/libs/utils"
)

type EnvConfig struct {
	ComposeDirPath string
	ServerPort     int
	SSHPrivateKey  string
	BinariesPath   string
	GitToken       string

	S3AccessKey string
	S3SecretKey string
	S3Endpoint  string
	S3Bucket    string
}

func getEnv() (env *EnvConfig) {
	utils.InitEnvFromFile()

	env = &EnvConfig{
		ComposeDirPath: utils.GetEnv("ORCHESTRATOR_COMPOSE_DIRECTORY_PATH", "/workspaces/home-config/compose"),
		ServerPort:     utils.GetEnv("ORCHESTRATOR_SERVER_PORT", 3000),
		SSHPrivateKey:  utils.GetEnv("ORCHESTRATOR_SSH_PRIVATE_KEY", ""),
		BinariesPath:   utils.GetEnv("ORCHESTRATOR_BINARIES_PATH", "/dist"),
		GitToken:       utils.GetEnv("ORCHESTRATOR_GIT_TOKEN", ""),
		S3AccessKey:    utils.GetEnv("ORCHESTRATOR_S3_ACCESS_KEY", ""),
		S3SecretKey:    utils.GetEnv("ORCHESTRATOR_S3_SECRET_KEY", ""),
		S3Endpoint:     utils.GetEnv("ORCHESTRATOR_S3_ENDPOINT", ""),
		S3Bucket:       utils.GetEnv("ORCHESTRATOR_S3_BUCKET", ""),
	}

	println("Using compose directory:", env.ComposeDirPath)

	return env
}

var Env *EnvConfig = getEnv()
