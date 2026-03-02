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
	ApiUrl     string
	ComposeDir string
	Mode       Mode
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

	// show os.Args
	fmt.Printf("os.Args: %v\n", os.Args)

	fmt.Printf("Mode : %s\n", mode)

	if mode == "" {
		mode = utils.GetEnv("ORCHESTRATOR_MODE", string(ModeHybrid))
	}

	env = &EnvConfig{
		ApiUrl:     utils.GetEnv("ORCHESTRATOR_API_URL", "http://localhost:3000"),
		ComposeDir: utils.GetEnv("ORCHESTRATOR_COMPOSE_DIRECTORY", ""),
		Mode:       Mode(mode),
	}

	return env
}

var Env *EnvConfig = getEnv()
