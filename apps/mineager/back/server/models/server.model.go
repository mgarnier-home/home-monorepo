package models

import (
	"context"
	"fmt"
	"mgarnier11/go/logger"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type ServerBo struct {
	Id      string
	Name    string
	Version string
	Map     string
	Url     string
	Memory  string
	Port    uint16
}

func mapContainerToServer(container *types.Container) *ServerBo {
	return &ServerBo{
		Id:      container.ID,
		Name:    container.Labels["mineager.serverName"],
		Version: container.Labels["mineager.serverVersion"],
		Map:     container.Labels["mineager.serverMap"],
		Url:     container.Labels["mineager.serverUrl"],
		Memory:  container.Labels["mineager.serverMemory"],
		// Port:    container.Ports[0].PublicPort,
	}
}

func getFilterArgs(name string) filters.Args {
	filterArgs := filters.NewArgs()
	// filterArgs.Add("label", "com.docker.compose.project=mineager")
	filterArgs.Add("ancestor", "itzg/minecraft-server")
	if name != "" {
		filterArgs.Add("label", "mineager.serverName="+name)
	}

	return filterArgs
}

func GetNextPort(dockerClient *client.Client, minPort uint16) (uint16, error) {
	servers, err := GetServers(dockerClient, "")

	if err != nil {
		return 0, err
	}

	var maxPort uint16 = minPort

	for _, server := range servers {
		if server.Port > maxPort {
			maxPort = server.Port
		}
	}

	return maxPort + 1, nil
}

func GetServers(dockerClient *client.Client, name string) (servers []*ServerBo, error error) {
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
		Filters: getFilterArgs(name),
		All:     true,
	})

	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		// logger.Infof("Container %v", container)

		server := mapContainerToServer(&container)

		servers = append(servers, server)
	}

	for _, server := range servers {
		logger.Infof("Server name: %s, version: %s, map: %s, url: %s, memory: %s, port: %d", server.Name, server.Version, server.Map, server.Url, server.Memory, server.Port)
	}

	return servers, nil
}

func CreateServer(dockerClient *client.Client, name string, version string, mapName string, memory string, url string, port uint16) (*ServerBo, error) {
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
		ExposedPorts: map[nat.Port]struct{}{
			"25565/tcp": {},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", port),
				},
			},
		},
	}

	ctx := context.Background()
	createResponse, err := dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, fmt.Sprintf("mineager-%s", name))
	if err != nil {
		return nil, err
	}

	err = dockerClient.ContainerStart(ctx, createResponse.ID, container.StartOptions{})
	if err != nil {
		return nil, err
	}

	newContainer, err := GetServers(dockerClient, name)
	if err != nil {
		return nil, err
	}

	return newContainer[0], nil

	// Connection en scp pour envoyer la map
	// Appel de l'api de infrared pour créer le reverse proxy
	// Création du container

}

func DeleteServer(dockerClient *client.Client, name string) error {
	servers, err := GetServers(dockerClient, name)

	if err != nil {
		return err
	}

	if len(servers) == 0 {
		return fmt.Errorf("server not found")
	}

	err = dockerClient.ContainerRemove(context.Background(), servers[0].Id, container.RemoveOptions{
		Force: true,
	})

	if err != nil {
		return err
	}

	return nil
}
