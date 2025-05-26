package config

import (
	"log"

	"mgarnier11.fr/go/libs/utils"
)

type Action struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}

type HealthCheck struct {
	Name   string  `yaml:"name"`
	Url    string  `yaml:"url"`
	Action *Action `yaml:"action"`
}

type Service struct {
	Name         string         `yaml:"name"`
	Icon         string         `yaml:"icon"`
	DockerName   string         `yaml:"dockerName"`
	HealthChecks []*HealthCheck `yaml:"healthChecks"`
}

type Host struct {
	Name      string     `yaml:"name"`
	Ip        string     `yaml:"ip"`
	Nodesight string     `yaml:"nodesight"`
	Icon      string     `yaml:"icon"`
	Services  []*Service `yaml:"services"`
}

type DashboardConfig struct {
	Hosts       []*Host `yaml:"hosts"`
	InfraredUrl string  `yaml:"infraredUrl"`
}

type AppEnvConfig struct {
	ServerPort     int
	AppDistPath    string
	ConfigFilePath string
	IconsPath      string

	DashboardConfig *DashboardConfig
}

func getAppEnvConfig() (appEnvConfig *AppEnvConfig) {
	utils.InitEnvFromFile()

	appEnvConfig = &AppEnvConfig{
		ServerPort:     utils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: utils.GetEnv("CONFIG_FILE_PATH", "./data/config.yaml"),
		AppDistPath:    utils.GetEnv("APP_DIST_PATH", "./dist"),
		IconsPath:      utils.GetEnv("ICONS_PATH", "./icons"),
	}

	var err error

	appEnvConfig.DashboardConfig, err = utils.ReadYamlFile[DashboardConfig](appEnvConfig.ConfigFilePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		panic(err)
	}

	return appEnvConfig
}

var Config *AppEnvConfig = getAppEnvConfig()
