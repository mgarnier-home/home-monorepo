package controllers

import (
	"context"
	"fmt"
	"mgarnier11/go/dockerssh"
	"mgarnier11/go/logger"
	"mgarnier11/go/utils"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/external"
	"mgarnier11/mineager/server/database"
	"mgarnier11/mineager/server/objects/bo"
	"mgarnier11/mineager/server/objects/dto"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func getDockerClient(host *config.DockerHostConfig) (*client.Client, error) {
	return dockerssh.GetDockerClient(host.SSHUsername, host.Ip, host.SSHPort, config.Config.SSHKeyPath)
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

func mapDockerContainerToServerBo(container *types.Container, inspect *types.ContainerJSON) *bo.ServerBo {
	port := uint16(0)
	containerPorts := inspect.HostConfig.PortBindings["25565/tcp"]
	if len(containerPorts) > 0 {
		intPort, err := strconv.Atoi(containerPorts[0].HostPort)
		if err == nil {
			port = uint16(intPort)
		}
	}

	return &bo.ServerBo{
		Id:      container.ID,
		Status:  inspect.State.Status,
		Name:    container.Labels["mineager.serverName"],
		Version: container.Labels["mineager.serverVersion"],
		Map:     container.Labels["mineager.serverMap"],
		Url:     container.Labels["mineager.serverUrl"],
		Memory:  container.Labels["mineager.serverMemory"],
		NewMap:  container.Labels["mineager.newMap"] == "true",
		Port:    port,
	}
}

type ServersController struct {
	host         *config.DockerHostConfig
	dockerClient *client.Client
	mapRepo      *database.MapRepository
}

func NewServersController(hostName string) (*ServersController, error) {
	host, err := config.GetHost(hostName)
	if err != nil {
		return nil, err
	}

	dockerClient, err := getDockerClient(host)
	if err != nil {
		return nil, err
	}

	return &ServersController{
		host:         host,
		dockerClient: dockerClient,
		mapRepo:      database.CreateMapRepository(),
	}, nil
}

func (controller *ServersController) Dispose() {
	controller.dockerClient.Close()
}

func (controller *ServersController) GetServers() ([]*bo.ServerBo, error) {
	containers, err := controller.dockerClient.ContainerList(context.Background(), container.ListOptions{
		Filters: getFilterArgs(""),
		All:     true,
	})

	if err != nil {
		return nil, err
	}

	servers := make([]*bo.ServerBo, 0)

	for _, container := range containers {
		containerInspect, err := controller.dockerClient.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return nil, err
		}

		servers = append(servers, mapDockerContainerToServerBo(&container, &containerInspect))
	}

	return servers, nil
}

func (controller *ServersController) GetServer(name string) (*bo.ServerBo, error) {
	servers, err := controller.GetServers()

	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.Name == name {
			return server, nil
		}
	}

	return nil, fmt.Errorf("server %s not found", name)
}

func (controller *ServersController) getNextPort() (uint16, error) {
	servers, err := controller.GetServers()

	if err != nil {
		return 0, err
	}

	var maxPort uint16 = uint16(controller.host.StartPort)

	for _, server := range servers {
		if server.Port > maxPort {
			maxPort = server.Port
		}
	}

	return maxPort + 1, nil
}

func (controller *ServersController) ServerExists(name string) (bool, error) {
	server, err := controller.GetServer(name)

	if server != nil {
		return true, nil
	} else if err != nil && err.Error() == fmt.Sprintf("server %s not found", name) {
		return false, nil
	} else {
		return false, err
	}
}

func (controller *ServersController) CreateServer(createServerDto *dto.CreateServerRequestDto) (*bo.ServerBo, error) {
	// Check if server already exists
	if serverExist, err := controller.ServerExists(createServerDto.Name); err != nil || serverExist {
		if err != nil {
			return nil, fmt.Errorf("error getting server: %v", err)
		}
		return nil, fmt.Errorf("server %s already exists", createServerDto.Name)
	}

	// Create server directory
	if err := createServerDirectory(createServerDto.Name); err != nil {
		return nil, fmt.Errorf("error creating server directory: %v", err)
	}

	// If the map is not new, copy it to the server directory and send it to the host
	if !createServerDto.NewMap {
		m, err := controller.mapRepo.GetMapByName(createServerDto.MapName)

		if err != nil {
			return nil, fmt.Errorf("error getting map: %v", err)
		}

		if m == nil {
			return nil, fmt.Errorf("map %s not found", createServerDto.MapName)
		}

		// Copy map to server directory
		if err := utils.CopyFolder(getMapPath(createServerDto.MapName), getServerLocalMapPath(createServerDto.Name)); err != nil {
			return nil, fmt.Errorf("error copying map to server directory: %v", err)
		}
	} else {
		_, err := controller.mapRepo.CreateMap(createServerDto.MapName, createServerDto.Version, "")
		if err != nil {
			return nil, fmt.Errorf("error creating a new  map: %v", err)
		}

		if err := os.MkdirAll(getServerLocalMapPath(createServerDto.Name), 0755); err != nil {
			return nil, fmt.Errorf("error creating new map directory: %v", err)
		}
	}

	// Send map to host
	// It will either send the copied map or a just create the dir on the host to allow the container to mount it and create a new map at start
	if err := sendServerMapToHost(createServerDto.Name, controller.host); err != nil {
		return nil, fmt.Errorf("error sending map to host: %v", err)
	}

	logger.Infof("Map %s sent to host %s on server %s", createServerDto.MapName, controller.host.Name, createServerDto.Name)
	port, err := controller.getNextPort()
	if err != nil {
		return nil, fmt.Errorf("error getting port: %v", err)
	}

	serverUrl := fmt.Sprintf("%s.%s", createServerDto.Name, config.Config.DomainName)

	containerConfig := &container.Config{
		Image: "itzg/minecraft-server", // Set la version du container
		Labels: map[string]string{
			"com.docker.compose.project": "mineager",
			"mineager.serverName":        createServerDto.Name,
			"mineager.serverVersion":     createServerDto.Version,
			"mineager.serverMap":         createServerDto.MapName,
			"mineager.serverUrl":         serverUrl,
			"mineager.serverMemory":      createServerDto.Memory,
			"mineager.newMap":            fmt.Sprintf("%t", createServerDto.NewMap),
			"traefik-conf.port":          fmt.Sprintf("%d", port),
		},
		Env: []string{
			"EULA=TRUE",
			"TYPE=VANILLA",
			"VERSION=" + createServerDto.Version,
			"MEMORY=" + createServerDto.Memory,
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
				Source: getServerHostMapPath(controller.host, createServerDto.Name),
				Target: "/data/world",
			},
		},
	}

	// Create the container
	createResponse, err := controller.dockerClient.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, fmt.Sprintf("mineager-%s", createServerDto.Name))
	if err != nil {
		return nil, fmt.Errorf("error creating container: %v", err)
	}

	serverBo := &bo.ServerBo{
		Id:      createResponse.ID,
		Name:    createServerDto.Name,
		Port:    port,
		Version: createServerDto.Version,
		Map:     createServerDto.MapName,
		Url:     serverUrl,
		Memory:  createServerDto.Memory,
	}

	if err := external.CreateProxy(controller.host, serverBo); err != nil {
		logger.Errorf("Error creating proxy: %v", err)
	}

	return serverBo, nil
}

func (controller *ServersController) StartServer(serverBo *bo.ServerBo) error {
	if err := sendServerMapToHost(serverBo.Name, controller.host); err != nil {
		return fmt.Errorf("error sending map to host: %v", err)
	}

	if err := controller.dockerClient.ContainerStart(context.Background(), serverBo.Id, container.StartOptions{}); err != nil {
		return fmt.Errorf("error starting container: %v", err)
	}

	return nil
}

func (controller *ServersController) StopServer(serverBo *bo.ServerBo) error {

	if err := controller.dockerClient.ContainerStop(context.Background(), serverBo.Id, container.StopOptions{}); err != nil {
		return fmt.Errorf("error stopping container: %v", err)
	}

	if err := getServerMapFromHost(serverBo.Name, controller.host); err != nil {
		return fmt.Errorf("error getting map from host: %v", err)
	}

	if serverBo.NewMap {
		if err := utils.CopyFolder(getServerLocalMapPath(serverBo.Name), getMapPath(serverBo.Map)); err != nil {
			return fmt.Errorf("error copying server map to maps directory: %v", err)
		}
	}

	return nil
}

func (controller *ServersController) DeleteServer(serverBo *bo.ServerBo, full bool) error {
	if err := external.DeleteProxy(controller.host, serverBo); err != nil {
		logger.Errorf("Error deleting proxy: %v", err)
	}

	if err := controller.dockerClient.ContainerRemove(context.Background(), serverBo.Id, container.RemoveOptions{
		Force: true,
	}); err != nil {
		return fmt.Errorf("error deleting container: %v", err)
	}

	if err := deleteHostDirectory(serverBo.Name, controller.host); err != nil {
		return fmt.Errorf("error deleting host server directory: %v", err)
	}

	if full {
		if err := deleteServerDirectory(serverBo.Name); err != nil {
			return fmt.Errorf("error deleting server directory: %v", err)
		}
	}

	return nil
}
