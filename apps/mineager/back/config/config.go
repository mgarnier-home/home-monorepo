package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"mgarnier11.fr/go/libs/utils"

	"gopkg.in/yaml.v3"
)

type AppConfigFile struct {
	InfraredUrl string              `yaml:"infraredUrl"`
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
	Name         string `yaml:"name"`
	Ip           string `yaml:"ip"`
	ProxyIp      string `yaml:"proxyIp"`
	SSHUsername  string `yaml:"sshUsername"`
	SSHPort      string `yaml:"sshPort"`
	StartPort    int    `yaml:"startPort"`
	MineagerPath string `yaml:"mineagerPath"`
}

type AppEnvConfig struct {
	ServerPort        int
	ConfigFilePath    string
	SSHKeyPath        string
	FrontendPath      string
	DataFolderPath    string
	MapsFolderPath    string
	ServersFolderPath string
	ApiToken          string
	DomainName        string

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
		DomainName:     utils.GetEnv("DOMAIN_NAME", ""),
	}

	appEnvConfig.MapsFolderPath = fmt.Sprintf("%s/maps", appEnvConfig.DataFolderPath)
	appEnvConfig.ServersFolderPath = fmt.Sprintf("%s/servers", appEnvConfig.DataFolderPath)

	err := os.MkdirAll(appEnvConfig.MapsFolderPath, 0755)
	if err != nil {
		log.Fatalf("Error creating maps folder: %v", err)
		panic(err)
	}

	err = os.MkdirAll(appEnvConfig.ServersFolderPath, 0755)
	if err != nil {
		log.Fatalf("Error creating servers folder: %v", err)
		panic(err)
	}

	appEnvConfig.AppConfig, err = readAppConfig(appEnvConfig.ConfigFilePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		panic(err)
	}

	return appEnvConfig
}

func GetHost(hostName string) (*DockerHostConfig, error) {
	for _, dockerHost := range Config.AppConfig.DockerHosts {
		if strings.EqualFold(dockerHost.Name, hostName) {
			return dockerHost, nil
		}
	}

	return nil, errors.New("host not found")
}

var Config *AppEnvConfig = getAppEnvConfig()
