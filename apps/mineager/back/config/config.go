package config

import (
	"fmt"
	"mgarnier11/go/logger"
	"mgarnier11/go/utils"
	"os"
)

type AppEnvConfig struct {
	ServerPort     int
	ConfigFilePath string
	SSHKeyPath     string
	FrontendPath   string
	DataFolderPath string
	MapsFolderPath string
	ApiToken       string
}

func getAppConfig() (appConfig *AppEnvConfig) {
	utils.InitEnvFromFile()

	appConfig = &AppEnvConfig{
		ServerPort:     utils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: utils.GetEnv("CONFIG_FILE_PATH", "./data/config.yaml"),
		SSHKeyPath:     utils.GetEnv("SSH_KEY_PATH", ""),
		FrontendPath:   utils.GetEnv("FRONTEND_PATH", "./front"),
		DataFolderPath: utils.GetEnv("DATA_FOLDER_PATH", "./data"),
		ApiToken:       utils.GetEnv("API_TOKEN", ""),
	}

	appConfig.MapsFolderPath = fmt.Sprintf("%s/maps", appConfig.DataFolderPath)

	err := os.MkdirAll(appConfig.MapsFolderPath, 0755)
	if err != nil {
		logger.Errorf("Error creating data and maps folder: %v", err)
		panic(err)
	}

	return appConfig
}

var Config *AppEnvConfig = getAppConfig()
