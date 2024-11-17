package hostmanager

import (
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/host"
	"strings"
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

func ConfigFileChanged(configFile *config.ConfigFile) {
	for _, hostConfig := range configFile.ProxyHosts {
		hostValue := GetHost(hostConfig.Name)

		if hostValue == nil {
			hostValue = host.NewHost(hostConfig)
			setHost(hostConfig.Name, hostValue)
		}
	}
}
