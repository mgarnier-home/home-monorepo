package hostManager

import (
	"slices"
	"strings"

	"mgarnier11.fr/go/libs/logger"

	"mgarnier11.fr/go/go-proxy/config"
	"mgarnier11.fr/go/go-proxy/host"
)

var hosts map[string]*host.Host = make(map[string]*host.Host)

func GetHosts() *(map[string]*host.Host) {
	return &hosts
}

func GetHost(name string) *host.Host {
	hostKey := strings.ToUpper(name)

	return hosts[hostKey]
}

func setHost(name string, host *host.Host) {
	hostKey := strings.ToUpper(name)

	hosts[hostKey] = host
}

func ConfigFileChanged(configFile *config.AppConfigFile) {
	logger.Infof("Config file changed")

	for hostKey, hostValue := range hosts {
		exists := slices.ContainsFunc(configFile.ProxyHosts, func(hostConfig *config.HostConfig) bool {
			return strings.ToUpper(hostConfig.Name) == hostKey
		})

		if !exists {
			logger.Infof("%s not found in updated config file, destroying it", hostValue.Config.Name)

			hostValue.Dispose()

			delete(hosts, hostKey)
		}
	}

	for _, hostConfig := range configFile.ProxyHosts {
		hostValue := GetHost(hostConfig.Name)

		if hostValue == nil {
			hostValue = host.NewHost(hostConfig)
			setHost(hostConfig.Name, hostValue)
		} else {
			hostValue.Config = hostConfig
		}
	}
}
