package config

import (
	"log"

	"mgarnier11.fr/go/libs/utils"
)

type AppConfigFile struct {
	InfraredUrl string `yaml:"infraredUrl"`
}

type AppEnvConfig struct {
	ServerPort     int
	ConfigFilePath string

	AppConfig *AppConfigFile
}

func getAppEnvConfig() (appEnvConfig *AppEnvConfig) {
	utils.InitEnvFromFile()

	appEnvConfig = &AppEnvConfig{
		ServerPort:     utils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: utils.GetEnv("CONFIG_FILE_PATH", "./data/config.yaml"),
	}

	var err error
	appEnvConfig.AppConfig, err = utils.ReadYamlFile[AppConfigFile](appEnvConfig.ConfigFilePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		panic(err)
	}

	return appEnvConfig
}

var Config *AppEnvConfig = getAppEnvConfig()
