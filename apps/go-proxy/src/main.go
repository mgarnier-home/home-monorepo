package main

import (
	"fmt"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/host"
	"runtime"
	"time"

	_ "net/http/pprof"

	"github.com/charmbracelet/log"
)

var hosts map[string]*host.Host = make(map[string]*host.Host)

func main() {
	log.SetLevel(log.DebugLevel)

	go func() {
		for {
			time.Sleep(time.Second)

			// Analyzing goroutine leaks
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)
			fmt.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())
		}
	}()

	for newConfigFile := range config.SetupConfigListener() {
		log.Infof("Config file changed")

		for _, host := range hosts {
			host.Dispose()

			log.Infof("Finished disposing %s", host.Config.Name)

		}

		hosts = make(map[string]*host.Host)

		for _, hostConfig := range newConfigFile.ProxyHosts {
			hosts[hostConfig.Name] = host.NewHost(hostConfig)
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
