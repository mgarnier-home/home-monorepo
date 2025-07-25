package config

import (
	"log"

	"mgarnier11.fr/go/libs/utils"
)

type AppConfigFile struct {
	KeepAliveUrl string            `yaml:"keepAliveUrl"`
	Mail         *MailConfig       `yaml:"mail"`
	FileName     string            `yaml:"fileName"`
	BackupSrc    string            `yaml:"backupSrc"`
	LocalDest    string            `yaml:"localDest"`
	RemoteDest   *RemoteDestConfig `yaml:"remoteDest"`
}

type MailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	InfoTo   string `yaml:"infoTo"`
	ErrorTo  string `yaml:"errorTo"`
}

type RemoteDestConfig struct {
	SSHHost string `yaml:"sshHost"`
	SSHPort int    `yaml:"sshPort"`
	SSHUser string `yaml:"sshUser"`
	SSHPath string `yaml:"sshPath"`
}

type AppEnvConfig struct {
	ServerPort     int
	ConfigFilePath string
	SSHPrivateKey  string

	AppConfig *AppConfigFile
}

func getAppEnvConfig() (appEnvConfig *AppEnvConfig) {
	utils.InitEnvFromFile()

	appEnvConfig = &AppEnvConfig{
		ServerPort:     utils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: utils.GetEnv("CONFIG_FILE_PATH", "./data/config.yaml"),
		SSHPrivateKey:  utils.GetEnv("SSH_PRIVATE_KEY", ""),
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
