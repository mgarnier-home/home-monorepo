package config

import (
	"fmt"
	"log"
	"mgarnier11/go/utils"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfigFile struct {
	DockerHosts []*DockerHostConfig `yaml:"dockerHosts"`
}

func readAppConfig(filePath string) (*AppConfigFile, error) {
	configFile, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	appConfig := &AppConfigFile{}

	err = yaml.Unmarshal(configFile, appConfig)

	if err != nil {
		return nil, err
	}

	return appConfig, nil
}

type DockerHostConfig struct {
	Name        string `yaml:"name"`
	Ip          string `yaml:"ip"`
	SSHUsername string `yaml:"sshUsername"`
	SSHPort     string `yaml:"sshPort"`
}

type AppEnvConfig struct {
	ServerPort     int
	ConfigFilePath string
	SSHKeyPath     string
	FrontendPath   string
	DataFolderPath string
	MapsFolderPath string
	ApiToken       string

	AppConfig *AppConfigFile
}

func getAppEnvConfig() (appEnvConfig *AppEnvConfig) {
	utils.InitEnvFromFile()

	appEnvConfig = &AppEnvConfig{
		ServerPort:     utils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: utils.GetEnv("CONFIG_FILE_PATH", "./data/config.yaml"),
		SSHKeyPath:     utils.GetEnv("SSH_KEY_PATH", ""),
		FrontendPath:   utils.GetEnv("FRONTEND_PATH", "./front"),
		DataFolderPath: utils.GetEnv("DATA_FOLDER_PATH", "./data"),
		ApiToken:       utils.GetEnv("API_TOKEN", ""),
	}

	appEnvConfig.MapsFolderPath = fmt.Sprintf("%s/maps", appEnvConfig.DataFolderPath)

	err := os.MkdirAll(appEnvConfig.MapsFolderPath, 0755)
	if err != nil {
		log.Fatalf("Error creating maps folder: %v", err)
		panic(err)
	}

	appEnvConfig.AppConfig, err = readAppConfig(appEnvConfig.ConfigFilePath)

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		panic(err)
	}

	return appEnvConfig
}

var Config *AppEnvConfig = getAppEnvConfig()
