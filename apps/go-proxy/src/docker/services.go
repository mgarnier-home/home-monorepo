package docker

import (
	"context"
	"errors"
	"fmt"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/cli/cli/connhelper"
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
		ServerPort: port,
		Protocol:   "tcp",
		Name:       containerName,
	}

	return proxyConfig, nil
}

func GetDockerClient(sshUsername string, hostIp string, sshPort int) (*client.Client, error) {
	helper, err := connhelper.GetConnectionHelper(fmt.Sprintf("ssh://%s@%s:%d", sshUsername, hostIp, sshPort))

	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
		Timeout: 2 * time.Second,
	}

	var clientOpts []client.Opt

	clientOpts = append(clientOpts,
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
	)

	client, err := client.NewClientWithOpts(clientOpts...)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetProxiesFromDocker(sshUsername string, hostIp string, logger *utils.Logger) ([]*config.ProxyConfig, error) {

	dockerClient, err := GetDockerClient(sshUsername, hostIp, 22)

	if err != nil {
		return nil, err
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		return nil, err
	}

	proxies := []*config.ProxyConfig{}

	for _, container := range containers {
		containerName := strings.Replace(container.Names[0], "/", "", 1)

		additionalPorts := container.Labels["proxy.ports"]

		proxyConfig, err := checkPortAndAddService(containerName, container.Labels["traefik-conf.port"])

		if err != nil {
			if logger != nil {
				logger.Verbosef("Error while checking port and adding service for container %s: %v", containerName, err)
			}
			continue
		}

		proxies = append(proxies, proxyConfig)

		if additionalPorts != "" {
			ports := strings.Split(additionalPorts, ",")

			for _, port := range ports {
				proxyConfig, err := checkPortAndAddService(containerName, port)

				if err != nil {
					if logger != nil {
						logger.Verbosef("Error adding additionnal ports for container %s: %v", containerName, err)
					}
					continue
				}

				proxies = append(proxies, proxyConfig)
			}
		}
	}

	return proxies, nil
}
