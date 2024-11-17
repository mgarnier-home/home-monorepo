package config

import (
	"os"
	"path"
	"time"

	"goUtils"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/log"
)

type ProxyConfig struct {
	ListenPort int    `yaml:"listenPort"`
	ServerPort int    `yaml:"serverPort"`
	Protocol   string `yaml:"protocol"`
	Name       string `yaml:"name"`
}

type HostConfig struct {
	Proxies      []*ProxyConfig `yaml:"proxies"`
	Name         string         `yaml:"name"`
	Ip           string         `yaml:"ip"`
	MacAddress   string         `yaml:"macAddress"`
	SSHUsername  string         `yaml:"sshUsername"`
	SSHPassword  string         `yaml:"sshPassword"`
	Autostop     bool           `yaml:"autostop"`
	MaxAliveTime int            `yaml:"maxAliveTime"`
}

type ConfigFile struct {
	ProxyHosts []*HostConfig `yaml:"proxyHosts"`
}

type AppConfig struct {
	ServerPort     int
	ConfigFilePath string
}

func readFile(filePath string) []byte {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	log.Infof("Read file %s", filePath)

	return bytes
}

func parseConfigFile(rawFile []byte) *ConfigFile {
	config := &ConfigFile{}
	err := yaml.Unmarshal(rawFile, config)
	if err != nil {
		panic(err)
	}

	return config
}
func GetAppConfig() (appConfig *AppConfig, err error) {
	envFilePath := goUtils.GetEnv("ENV_FILE_PATH", "./.env")

	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	exPath := path.Dir(ex)

	if path.IsAbs(envFilePath) {
		envFilePath = path.Join(exPath, envFilePath)
	}

	godotenv.Load(envFilePath)

	appConfig = &AppConfig{
		ServerPort:     goUtils.GetEnv("SERVER_PORT", 8080),
		ConfigFilePath: goUtils.GetEnv("CONFIG_FILE_PATH", "config.yaml"),
	}

	return appConfig, nil
}

func SetupConfigListener() chan *ConfigFile {
	newConfigFileChan := make(chan *ConfigFile)

	appConfig, err := GetAppConfig()

	if err != nil {
		panic(err)
	}

	go func() {
		// Read and send the initial config file
		oldYamlFile := readFile(appConfig.ConfigFilePath)
		newConfigFileChan <- parseConfigFile(oldYamlFile)
		for range time.Tick(time.Second * 5) {
			yamlFile := readFile(appConfig.ConfigFilePath)

			if string(yamlFile) != string(oldYamlFile) {
				newConfigFileChan <- parseConfigFile(yamlFile)
			}

			oldYamlFile = yamlFile
		}
	}()

	return newConfigFileChan
}
