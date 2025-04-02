package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/utils"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type ProxyConfig struct {
	ListenPort int    `yaml:"listenPort"`
	ServerPort int    `yaml:"serverPort"`
	Protocol   string `yaml:"protocol"`
	Name       string `yaml:"name"`
	Key        string
}

type HostConfig struct {
	Proxies      []*ProxyConfig `yaml:"proxies"`
	Name         string         `yaml:"name"`
	Ip           string         `yaml:"ip"`
	MacAddress   string         `yaml:"macAddress"`
	SSHUsername  string         `yaml:"sshUsername"`
	SSHPort      string         `yaml:"sshPort"`
	Autostop     bool           `yaml:"autostop"`
	MaxAliveTime int            `yaml:"maxAliveTime"`

	appConfig *AppConfigFile
}

func (hostConfig *HostConfig) Save() {
	hostConfig.appConfig.Save()
}

type AppConfigFile struct {
	ProxyHosts []*HostConfig `yaml:"proxyHosts"`
}

func (config *AppConfigFile) unmarshalYaml(yamlString string) error {
	err := yaml.Unmarshal([]byte(yamlString), config)
	if err != nil {
		return err
	}

	for _, hostConfig := range config.ProxyHosts {
		hostConfig.appConfig = config
	}

	return nil

}

func (config *AppConfigFile) Save() {
	err := saveConfigFile(config, Config.ConfigFilePath)
	if err != nil {
		logger.Errorf("Failed to save config file: %s", err)
	}
}

type AppEnvConfig struct {
	ServerPort     int
	ConfigFilePath string
	SSHKeyPath     string
}

func readFile(filePath string) []byte {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return bytes
}

func parseConfigFile(rawFile []byte) *AppConfigFile {
	config := &AppConfigFile{}

	err := config.unmarshalYaml(string(rawFile))

	if err != nil {
		panic(err)
	}

	for _, hostConfig := range config.ProxyHosts {
		for _, proxyConfig := range hostConfig.Proxies {
			proxyConfig.Key = fmt.Sprintf("%s:%d", proxyConfig.Name, proxyConfig.ListenPort)
		}
	}

	return config
}

func saveConfigFile(config *AppConfigFile, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	return os.WriteFile(filePath, data, 0644)
}

func getAppConfig() (appConfig *AppEnvConfig) {
	envFilePath := utils.GetEnv("ENV_FILE_PATH", "./.env")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := path.Dir(ex)

	if path.IsAbs(envFilePath) {
		envFilePath = path.Join(exPath, envFilePath)
	}

	godotenv.Load(envFilePath)

	appConfig = &AppEnvConfig{
		ServerPort:     utils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: utils.GetEnv("CONFIG_FILE_PATH", "config.yaml"),
		SSHKeyPath:     utils.GetEnv("SSH_KEY_PATH", ""),
	}

	return appConfig
}

func SetupConfigListener() chan *AppConfigFile {
	newConfigFileChan := make(chan *AppConfigFile)

	go func() {
		// Read and send the initial config file
		oldYamlFile := readFile(Config.ConfigFilePath)
		newConfigFileChan <- parseConfigFile(oldYamlFile)
		for range time.Tick(time.Second * 5) {
			yamlFile := readFile(Config.ConfigFilePath)

			if string(yamlFile) != string(oldYamlFile) {
				newConfigFileChan <- parseConfigFile(yamlFile)
			}

			oldYamlFile = yamlFile
		}
	}()

	return newConfigFileChan
}

var Config *AppEnvConfig = getAppConfig()
