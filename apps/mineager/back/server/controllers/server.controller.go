package controllers

import (
	"fmt"
	"mgarnier11/go/dockerssh"
	"mgarnier11/go/logger"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/server/models"

	"github.com/docker/docker/client"
)

type ServerDto struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Version string `json:"version"`
	Map     string `json:"map"`
}

func serverBoToServerDto(serverBo *models.ServerBo) *ServerDto {
	return &ServerDto{
		Name:    serverBo.Name,
		Url:     serverBo.Url,
		Version: serverBo.Version,
		Map:     serverBo.Map,
	}
}

func serversBoToServersDto(serversBo []*models.ServerBo) []*ServerDto {
	serversDto := make([]*ServerDto, 0)

	for _, serverBo := range serversBo {
		serversDto = append(serversDto, serverBoToServerDto(serverBo))
	}

	return serversDto
}

func getDockerClient(host *config.DockerHostConfig) (*client.Client, error) {
	return dockerssh.GetDockerClient(host.SSHUsername, host.Ip, host.SSHPort, config.Config.SSHKeyPath)
}

func getDockerClients(hostName string) ([]*client.Client, error) {
	clients := make([]*client.Client, 0)

	host, err := config.GetHost(hostName)

	if err != nil {
		for _, dockerHost := range config.Config.AppConfig.DockerHosts {
			dockerClient, err := getDockerClient(dockerHost)

			if err != nil {
				return nil, err
			}

			clients = append(clients, dockerClient)
		}
	} else {
		dockerClient, err := getDockerClient(host)

		if err != nil {
			return nil, err
		}

		clients = append(clients, dockerClient)
	}

	return clients, nil
}

func GetServers(hostName string, name string) ([]*ServerDto, error) {
	servers := make([]*ServerDto, 0)

	dockerClients, err := getDockerClients(hostName)

	if err != nil {
		return nil, err
	}

	for _, dockerClient := range dockerClients {
		defer dockerClient.Close()

		serversBo, err := models.GetServers(dockerClient, name)

		if err != nil {
			logger.Errorf("error getting servers %v", err)
			continue
		}

		servers = append(servers, serversBoToServersDto(serversBo)...)
	}

	return servers, nil
}

func CreateServer(
	hostName string,
	serverName string,
	version string,
	mapName string,
	memory string,
	url string,
) (*ServerDto, error) {
	host, err := config.GetHost(hostName)

	if err != nil {
		return nil, err
	}

	dockerClient, err := getDockerClient(host)
	if err != nil {
		return nil, err
	}
	defer dockerClient.Close()

	serverExist, err := models.ServerExists(dockerClient, serverName)
	if err != nil {
		return nil, fmt.Errorf("error getting server: %v", err)
	}
	if serverExist {
		return nil, fmt.Errorf("server %s already exists", serverName)
	}

	// Copy map to servers folder
	mapPath := getMapPath(mapName)

	// err = sendMapToHost(serverName, mapName, host)
	// if err != nil {
	// 	return nil, fmt.Errorf("error sending map to host: %v", err)
	// }
	// logger.Infof("Map %s sent to host %s on server %s", mapName, hostName, serverName)

	serverConfig := &models.ServerConfig{
		Client:  dockerClient,
		Host:    host,
		Name:    serverName,
		Version: version,
		Map:     mapName,
		Memory:  memory,
		Url:     url,
	}
	serverBo, err := models.CreateServer(serverConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating server: %v", err)
	}

	return serverBoToServerDto(serverBo), nil
}

func DeleteServer(hostName string, serverName string) error {
	host, err := config.GetHost(hostName)

	if err != nil {
		return err
	}

	dockerClient, err := getDockerClient(host)

	if err != nil {
		return err
	}
	defer dockerClient.Close()

	return models.DeleteServer(dockerClient, serverName)
}
