package main

import (
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/hostmanager"
	"mgarnier11/go-proxy/server"
	"runtime"
	"time"

	_ "net/http/pprof"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.InfoLevel)

	go func() {
		for {
			time.Sleep(time.Second * 5)

			// Analyzing goroutine leaks
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)
			log.Infof("Number of Goroutines: %d", runtime.NumGoroutine())

			// buf := make([]byte, 1<<16) // Create a large buffer to capture stack traces
			// stackLen := runtime.Stack(buf, true)
			// fmt.Printf("=== Goroutine Stack Dump ===\n%s\n", buf[:stackLen])
		}
	}()

	appConfig, err := config.GetAppConfig()

	if err != nil {
		panic(err)
	}

	server := server.NewServer(appConfig.ServerPort)

	go server.Start()

	for configFile := range config.SetupConfigListener() {
		hostmanager.ConfigFileChanged(configFile)
	}

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
