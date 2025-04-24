package config

import (
	"log"

	"mgarnier11.fr/go/libs/utils"
)

type AppConfigFile struct {
	Mail       MailConfig       `yaml:"mail"`
	BackupDest BackupDestConfig `yaml:"backupDest"`
}

type MailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	InfoTo   string `yaml:"infoTo"`
	ErrorTo  string `yaml:"errorTo"`
}

type BackupDestConfig struct {
	SSHHost string `yaml:"sshHost"`
	SSHPort int    `yaml:"sshPort"`
	SSHUser string `yaml:"sshUser"`
	SSHKey  string `yaml:"sshKey"`
	SSHPath string `yaml:"sshPath"`
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
