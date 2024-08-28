package main

import (
	"context"
	"log"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/docker"
)

var hostConfig = &config.HostConfig{
	Name:       "test",
	Ip:         "100.64.98.100",
	DockerPort: 4321,
}

var proxyConfig = &config.ProxyConfig{
	ListenPort: 8080,
	TargetPort: 12004,
	Protocol:   "tcp",
	Name:       "test",
}

func main() {
	println("Hello, World!")
	context := context.Background()

	proxies, err := docker.GetProxiesFromDocker(context, hostConfig.Ip, hostConfig.DockerPort)

	if err != nil {
		panic(err)
	}

	for _, proxy := range proxies {
		log.Printf("Proxy: %v\n", proxy)
	}

	// tcpProxy := proxies.NewTCPProxy(context, &proxies.TCPProxyArgs{
	// 	ProxyConfig: proxyConfig,
	// 	HostConfig:  hostConfig,
	// 	HostStarted: func(proxy *proxies.TCPProxy) (bool, error) {
	// 		return true, nil
	// 	},
	// 	StartHost: func(proxy *proxies.TCPProxy) error {
	// 		return nil
	// 	},
	// 	PacketReceived: func(proxy *proxies.TCPProxy) error {
	// 		return nil
	// 	},
	// })

	// var wg sync.WaitGroup
	// wg.Add(1)

	// tcpProxy.Start(&wg)

	// wg.Wait()
}
