package docker

import (
	"context"
	"errors"
	"fmt"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go/dockerssh"
	"mgarnier11/go/logger"
	sshUtils "mgarnier11/go/utils/ssh"

	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
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
		Key:        fmt.Sprintf("%s:%d", containerName, port),
	}

	return proxyConfig, nil
}

func GetDockerClient(sshUsername string, hostIp string, sshPort int) (*client.Client, error) {
	authMethod, err := sshUtils.GetSSHKeyAuth(config.Config.SSHKeyPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get ssh key auth: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            sshUsername,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Replace with a proper callback in production
	}

	sshDialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dockerssh.NewSSHDialer(
			net.JoinHostPort(hostIp, strconv.Itoa(sshPort)),
			sshConfig,
		)
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: sshDialer,
		},
		Timeout: 2 * time.Second,
	}

	var clientOpts []client.Opt

	clientOpts = append(clientOpts,
		client.WithHTTPClient(httpClient),
		client.WithDialContext(sshDialer),
	)

	client, err := client.NewClientWithOpts(clientOpts...)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetProxiesFromDocker(sshUsername string, hostIp string, logger *logger.Logger) ([]*config.ProxyConfig, error) {

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
