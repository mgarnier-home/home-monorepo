package models

import (
	"context"
	"fmt"
	"mgarnier11/go/dockerssh"
	"mgarnier11/go/logger"
	"mgarnier11/mineager/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

type ServerBo struct {
	Name    string
	Version string
	Map     string
	Url     string
	Memory  string
}

func getServersOnHost(dockerHost *config.DockerHostConfig) (servers []*ServerBo, error error) {

	logger.Infof("Getting servers on host %s", dockerHost.Name)
	logger.Infof("connecting to %s:%s", dockerHost.Ip, dockerHost.SSHPort)
	dockerClient, err := dockerssh.GetDockerClient(dockerHost.SSHUsername, dockerHost.Ip, dockerHost.SSHPort, config.Config.SSHKeyPath)

	if err != nil {
		return nil, err
	}
	defer dockerClient.Close()

	filterArgs := filters.NewArgs()
	// filterArgs.Add("label", "com.docker.compose.project=mineager")
	filterArgs.Add("ancestor", "itzg/minecraft-server")

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
		Filters: filterArgs,
	})

	if err != nil {
		return nil, err
	}

	logger.Infof("Containers : %v", containers)

	for _, container := range containers {
		server := &ServerBo{
			Name:    container.Labels["mineager.serverName"],
			Version: container.Labels["mineager.serverVersion"],
			Map:     container.Labels["mineager.serverMap"],
			Url:     container.Labels["mineager.serverUrl"],
			Memory:  container.Labels["mineager.serverMemory"],
		}

		servers = append(servers, server)
	}

	return servers, nil
}

func GetServers() (servers []*ServerBo, error error) {

	for _, dockerHost := range config.Config.AppConfig.DockerHosts {
		serversOnHost, err := getServersOnHost(dockerHost)

		if err != nil {
			logger.Errorf("error getting servers on host %s: %v", dockerHost.Name, err)
		} else {
			servers = append(servers, serversOnHost...)
		}
	}

	return servers, nil
}

func GetServerByName(name string) (*ServerBo, error) {
	for _, dockerHost := range config.Config.AppConfig.DockerHosts {
		serversOnHost, err := getServersOnHost(dockerHost)

		if err != nil {
			logger.Errorf("error getting servers on host %s: %v", dockerHost.Name, err)
		} else {
			for _, server := range serversOnHost {
				if server.Name == name {
					return server, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("server not found")
}

func createServerOnHost(dockerHost *config.DockerHostConfig, name string, version string, mapName string, url string, memory string) (*ServerBo, error) {

	logger.Infof("Creating server %s on host %s", name, dockerHost.Name)

	dockerClient, err := dockerssh.GetDockerClient(dockerHost.SSHUsername, dockerHost.Ip, dockerHost.SSHPort, config.Config.SSHKeyPath)

	if err != nil {
		return nil, err
	}
	defer dockerClient.Close()

	containerConfig := &container.Config{
		Image: "itzg/minecraft-server", // Set la version du container
		Labels: map[string]string{
			"mineager.serverName":    name,
			"mineager.serverVersion": version,
			"mineager.serverMap":     mapName,
			"mineager.serverUrl":     url,
			"mineager.serverMemory":  memory,
			// Label traefik conf
		},
		Env: []string{
			"EULA=TRUE",
			"TYPE=VANILLA",
			"VERSION=" + version,
			"MEMORY=" + memory,
		},
		// PORTS

	}

	// Connection en scp pour envoyer la map
	// Appel de l'api de infrared pour créer le reverse proxy
	// Création du container

}

func CreateServer(host string, name string, version string, mapName string, url string, memory string) (*ServerBo, error) {

	logger.Infof("Creating server %s on host %s", name, host)

	for _, dockerHost := range config.Config.AppConfig.DockerHosts {
		if dockerHost.Name == host {
			return createServerOnHost(dockerHost, name, version, mapName, url, memory)
		}
	}

	return nil, fmt.Errorf("host not found")
}
