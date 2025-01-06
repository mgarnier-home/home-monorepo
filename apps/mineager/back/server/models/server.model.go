package models

import (
	"context"
	"fmt"
	"mgarnier11/mineager/config"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
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

type ServerConfig struct {
	Name    string
	Version string
	Map     string
	Memory  string
	Url     string
	Host    *config.DockerHostConfig
	Client  *client.Client
}

func mapContainerToServer(container *types.Container) *ServerBo {
	port := uint16(0)
	if len(container.Ports) > 0 {
		port = container.Ports[0].PublicPort
	}

	return &ServerBo{
		Id:      container.ID,
		Name:    container.Labels["mineager.serverName"],
		Version: container.Labels["mineager.serverVersion"],
		Map:     container.Labels["mineager.serverMap"],
		Url:     container.Labels["mineager.serverUrl"],
		Memory:  container.Labels["mineager.serverMemory"],
		Port:    port,
	}
}

func getFilterArgs(name string) filters.Args {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "com.docker.compose.project=mineager")
	filterArgs.Add("ancestor", "itzg/minecraft-server")
	if name != "" {
		filterArgs.Add("label", "mineager.serverName="+name)
	}

	return filterArgs
}

func GetServerDestPath(host *config.DockerHostConfig, serverName string) string {
	return fmt.Sprintf("%s/%s/world", host.MineagerPath, serverName)
}

func getNextPort(dockerClient *client.Client, minPort uint16) (uint16, error) {
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

func ServerExists(dockerClient *client.Client, name string) (bool, error) {
	servers, err := GetServers(dockerClient, name)

	if err != nil {
		return false, err
	}

	return len(servers) > 0, nil
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
		server := mapContainerToServer(&container)

		servers = append(servers, server)
	}

	return servers, nil
}

func CreateServer(serverConfig *ServerConfig) (*ServerBo, error) {
	port, err := getNextPort(serverConfig.Client, uint16(serverConfig.Host.StartPort))

	if err != nil {
		return nil, fmt.Errorf("error getting port: %v", err)
	}

	containerConfig := &container.Config{
		Image: "itzg/minecraft-server", // Set la version du container
		Labels: map[string]string{
			"com.docker.compose.project": "mineager",
			"mineager.serverName":        serverConfig.Name,
			"mineager.serverVersion":     serverConfig.Version,
			"mineager.serverMap":         serverConfig.Map,
			"mineager.serverUrl":         serverConfig.Url,
			"mineager.serverMemory":      serverConfig.Memory,
			// Label traefik conf
		},
		Env: []string{
			"EULA=TRUE",
			"TYPE=VANILLA",
			"VERSION=" + serverConfig.Version,
			"MEMORY=" + serverConfig.Memory,
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
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: GetServerDestPath(serverConfig.Host, serverConfig.Name),
				Target: "/data/world",
			},
		},
	}

	createResponse, err := serverConfig.Client.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		nil,
		nil,
		fmt.Sprintf("mineager-%s", serverConfig.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating container: %v", err)
	}

	return &ServerBo{
		Id:      createResponse.ID,
		Name:    serverConfig.Name,
		Port:    port,
		Version: serverConfig.Version,
		Map:     serverConfig.Map,
		Url:     serverConfig.Url,
		Memory:  serverConfig.Memory,
	}, nil
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
