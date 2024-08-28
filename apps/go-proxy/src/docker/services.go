package docker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mgarnier11/go-proxy/config"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func checkPortAndAddService(containerName string, traefikConfPort string) (*config.ProxyConfig, error) {

	if traefikConfPort == "" {
		return nil, errors.New("traefik-conf.port not found")
	}

	port, err := strconv.Atoi(traefikConfPort)

	if err != nil {
		return nil, err
	}

	proxyConfig := &config.ProxyConfig{
		ListenPort: port,
		TargetPort: port,
		Protocol:   "tcp",
		Name:       containerName,
	}

	return proxyConfig, nil
}

func GetProxiesFromDocker(ctx context.Context, hostIp string, hostDockerPort int) ([]*config.ProxyConfig, error) {

	dockerClient, err := client.NewClientWithOpts(client.WithHost(fmt.Sprintf("tcp://%s:%d", hostIp, hostDockerPort)), client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{})

	if err != nil {
		panic(err)
	}

	proxies := []*config.ProxyConfig{}

	for _, container := range containers {
		containerName := strings.Replace(container.Names[0], "/", "", 1)

		log.Println(containerName)

		traefikConfPort := container.Labels["traefik-conf.port"]
		additionalPorts := container.Labels["proxy.ports"]

		proxyConfig, err := checkPortAndAddService(containerName, traefikConfPort)

		if err != nil {
			log.Println(err)
			continue
		}

		proxies = append(proxies, proxyConfig)

		if additionalPorts != "" {
			ports := strings.Split(additionalPorts, ",")

			for _, port := range ports {
				proxyConfig, err := checkPortAndAddService(containerName, port)

				if err != nil {
					log.Println(err)
					continue
				}

				proxies = append(proxies, proxyConfig)
			}
		}

	}

	return proxies, nil
}
