package main

import (
	"context"
	"log"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/host"
)

var hosts map[string]*host.Host = make(map[string]*host.Host)

func main() {
	context := context.Background()

	for newConfigFile := range config.SetupConfigListener() {
		log.Println("Config file changed")

		for _, host := range hosts {
			host.Dispose()

			log.Printf("Host %s disposed\n", host.Config.Name)

		}

		hosts = make(map[string]*host.Host)

		for _, hostConfig := range newConfigFile.ProxyHosts {
			hosts[hostConfig.Name] = host.NewHost(context, hostConfig)
		}

	}
	// appConfig, configFile, err := config.GetAppConfig()

	// if err != nil {
	// 	panic(err)
	// }

	// log.Printf("AppConfig: %+v\n", appConfig)
	// log.Printf("ConfigFile: %+v\n", configFile)

	// context := context.Background()

	// for _, hostConfig := range configFile.ProxyHosts {
	// 	host := host.NewHost(context, hostConfig)
	// }

	// context := context.Background()

	// proxies, err := docker.GetProxiesFromDocker(context, hostConfig.Ip, hostConfig.DockerPort)

	// if err != nil {
	// 	panic(err)
	// }

	// for _, proxy := range proxies {
	// 	log.Printf("Proxy: %v\n", proxy)
	// }

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
