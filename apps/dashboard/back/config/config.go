package config

import (
	"encoding/json"

	"mgarnier11.fr/go/libs/utils"
)

type Action struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
	Url  string `yaml:"url" json:"url"`
}

type HealthCheck struct {
	Name   string  `yaml:"name" json:"name"`
	Url    string  `yaml:"url" json:"url"`
	Action *Action `yaml:"action" json:"action"`
}

type Service struct {
	Name         string         `yaml:"name" json:"name"`
	Icon         string         `yaml:"icon" json:"icon"`
	DockerName   string         `yaml:"dockerName" json:"dockerName"`
	HealthChecks []*HealthCheck `yaml:"healthChecks" json:"healthChecks"`
}

type Host struct {
	Name      string     `yaml:"name" json:"name"`
	Ip        string     `yaml:"ip" json:"ip"`
	Nodesight string     `yaml:"nodesight" json:"nodesight"`
	Icon      string     `yaml:"icon" json:"icon"`
	Services  []*Service `yaml:"services" json:"services"`
}

type DashboardConfig struct {
	Hosts []*Host `yaml:"hosts" json:"hosts"`
}

func (dashboardConfig *DashboardConfig) ToJSON() (string, error) {
	bytes, err := json.Marshal(dashboardConfig)

	if err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

type AppEnvConfig struct {
	ServerPort     int
	AppDistPath    string
	ConfigFilePath string
	IconsPath      string
}

func (appEnvConfig *AppEnvConfig) GetDashboardConfig() (*DashboardConfig, error) {
	return utils.ReadYamlFile[DashboardConfig](appEnvConfig.ConfigFilePath)
}

func (appEnvConfig *AppEnvConfig) GetDashboardConfigJSON() (string, error) {
	dashboardConfig, err := appEnvConfig.GetDashboardConfig()
	if err != nil {
		return "", err
	} else {
		return dashboardConfig.ToJSON()
	}
}

func getAppEnvConfig() (appEnvConfig *AppEnvConfig) {
	utils.InitEnvFromFile()

	appEnvConfig = &AppEnvConfig{
		ServerPort:     utils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: utils.GetEnv("CONFIG_FILE_PATH", "./data/config.yaml"),
		AppDistPath:    utils.GetEnv("APP_DIST_PATH", "./dist"),
		IconsPath:      utils.GetEnv("ICONS_PATH", "./icons"),
	}

	return appEnvConfig
}

var Config *AppEnvConfig = getAppEnvConfig()
