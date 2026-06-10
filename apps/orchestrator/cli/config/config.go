package config

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
	"mgarnier11.fr/go/libs/utils"
)

type Mode string

const (
	ModeFullLocal Mode = "local"
	ModeFullApi   Mode = "remote"
	ModeHybrid    Mode = "hybrid"
)

type EnvConfig struct {
	ApiUrl         string
	ComposeDirPath string
	Mode           Mode
	S3AccessKey    string
	S3SecretKey    string
	S3Endpoint     string
	S3Bucket       string
}

func getEnv() (env *EnvConfig) {
	utils.InitEnvFromFile()

	var mode, service string

	bootstrap := pflag.NewFlagSet("bootstrap", pflag.ContinueOnError)
	bootstrap.SetOutput(io.Discard) // silence pre-parse usage/errors
	bootstrap.StringVar(&mode, "mode", "", "Set the execution mode")
	bootstrap.StringVar(&service, "service", "", "Execute command for a specific service")
	err := bootstrap.Parse(os.Args[1:]) // mode is now available here

	if err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
	}

	if mode == "" {
		mode = utils.GetEnv("ORCHESTRATOR_MODE", string(ModeHybrid))
	}

	env = &EnvConfig{
		ApiUrl:         utils.GetEnv("ORCHESTRATOR_API_URL", "http://localhost:3000"),
		ComposeDirPath: utils.GetEnv("ORCHESTRATOR_COMPOSE_DIRECTORY_PATH", "/workspaces/home-config/compose"),
		Mode:           Mode(mode),
		S3AccessKey:    utils.GetEnv("ORCHESTRATOR_S3_ACCESS_KEY", ""),
		S3SecretKey:    utils.GetEnv("ORCHESTRATOR_S3_SECRET_KEY", ""),
		S3Endpoint:     utils.GetEnv("ORCHESTRATOR_S3_ENDPOINT", ""),
		S3Bucket:       utils.GetEnv("ORCHESTRATOR_S3_BUCKET", ""),
	}

	return env
}

var Env *EnvConfig = getEnv()
