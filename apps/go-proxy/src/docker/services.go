package docker

import (
	"context"
	"errors"
	"fmt"
	"mgarnier11/go-proxy/config"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func SetupDockerContainersListener(ctx context.Context, hostIp string, hostDockerPort int) chan []*config.ProxyConfig {
	proxiesChan := make(chan []*config.ProxyConfig)

	go func() {
		getProxies := func() {
			proxies, err := GetProxiesFromDocker(hostIp, hostDockerPort)

			if err != nil {
				log.Errorf("Error while getting proxies from docker: %v", err)
			} else {
				proxiesChan <- proxies
			}
		}

		getProxies()

		ticker := time.NewTicker(5 * time.Second)

		go func() {
			log.Infof("Waiting for ctx.Done()")
			<-ctx.Done()
			log.Infof("ctx.Done() received")
			ticker.Stop()
			close(proxiesChan)
			log.Infof("Docker containers listener stopped")
		}()

		for range ticker.C {
			log.Infof("Getting proxies from docker")
			getProxies()
		}
	}()

	return proxiesChan
}

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

func GetProxiesFromDocker(hostIp string, hostDockerPort int) ([]*config.ProxyConfig, error) {

	dockerClient, err := client.NewClientWithOpts(client.WithHost(fmt.Sprintf("tcp://%s:%d", hostIp, hostDockerPort)), client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		panic(err)
	}

	proxies := []*config.ProxyConfig{}

	for _, container := range containers {
		containerName := strings.Replace(container.Names[0], "/", "", 1)

		traefikConfPort := container.Labels["traefik-conf.port"]
		additionalPorts := container.Labels["proxy.ports"]

		proxyConfig, err := checkPortAndAddService(containerName, traefikConfPort)

		if err != nil {
			log.Debugf("Error while checking port and adding service for container %s: %v", containerName, err)
			continue
		}

		proxies = append(proxies, proxyConfig)

		if additionalPorts != "" {
			ports := strings.Split(additionalPorts, ",")

			for _, port := range ports {
				proxyConfig, err := checkPortAndAddService(containerName, port)

				if err != nil {
					log.Debugf("Error while checking port and adding service for container %s: %v", containerName, err)
					continue
				}

				proxies = append(proxies, proxyConfig)
			}
		}

	}

	return proxies, nil
}
