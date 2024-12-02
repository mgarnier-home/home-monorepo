package docker

import (
	"context"
	"errors"
	"fmt"
	"mgarnier11/go-proxy/config"
	myconnhelper "mgarnier11/go/docker"
	"mgarnier11/go/logger"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	// // Load the private key file
	// key, err := os.ReadFile(config.Config.SSHKeyPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read private key: %w", err)
	// }

	// // Parse the private key
	// signer, err := ssh.ParsePrivateKey(key)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to parse private key: %w", err)
	// }

	// Configure the SSH client
	// sshConfig := &ssh.ClientConfig{
	// 	User: sshUsername,
	// 	// Auth: []ssh.AuthMethod{
	// 	// 	ssh.PublicKeys(signer),
	// 	// },
	// 	Auth:            []ssh.AuthMethod{ssh.Password("P@55w0rd")},
	// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Replace with a proper callback in production
	// }

	helper, err := myconnhelper.GetConnectionHelper(
		fmt.Sprintf("ssh://%s@%s:%d",
			sshUsername,
			hostIp,
			sshPort,
		),
	)

	// Custom DialContext using Go's SSH library
	// helper.Dialer = func(ctx context.Context, network, address string) (net.Conn, error) {
	// 	conn, err := ssh.Dial("tcp", net.JoinHostPort(hostIp, "22"), sshConfig)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	// 	}
	// 	return conn.Dial(network, address)
	// }

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
		// client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
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
