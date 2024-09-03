package host

import (
	"context"
	"log"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/proxies"
	"sync"
	"time"
)

type Host struct {
	TCPProxies     map[string]*proxies.TCPProxy
	UDPProxies     map[string]*proxies.UDPProxy
	Started        bool
	LastPacketDate time.Time
	Config         *config.HostConfig
	waitGroup      sync.WaitGroup
}

func NewHost(ctx context.Context, hostConfig *config.HostConfig) *Host {
	host := &Host{
		Config:     hostConfig,
		TCPProxies: make(map[string]*proxies.TCPProxy),
		UDPProxies: make(map[string]*proxies.UDPProxy),
		waitGroup:  sync.WaitGroup{},
	}

	for _, proxyConfig := range hostConfig.Proxies {
		if proxyConfig.Protocol == "tcp" {
			host.TCPProxies[proxyConfig.Name] = proxies.NewTCPProxy(context.Background(), &proxies.TCPProxyArgs{
				HostConfig:     hostConfig,
				ProxyConfig:    proxyConfig,
				HostStarted:    host.HostStarted,
				StartHost:      host.StartHost,
				PacketReceived: host.PacketReceived,
			})

			host.waitGroup.Add(1)
			go host.TCPProxies[proxyConfig.Name].Start(&host.waitGroup)

		} else if proxyConfig.Protocol == "udp" {

		}
	}

	host.StartHost(nil)
	host.Started = true
	host.LastPacketDate = time.Now()

	log.Println("Host created : " + host.Config.Name)

	return host
}

func (host *Host) HostStarted(proxy *proxies.TCPProxy) (bool, error) {
	return host.Started, nil
}

func (host *Host) StartHost(proxy *proxies.TCPProxy) error {
	log.Println("Starting host : " + host.Config.Name)
	return nil
}

func (host *Host) StopHost() error {
	log.Println("Stopping host : " + host.Config.Name)
	return nil
}

func (host *Host) PacketReceived(proxy *proxies.TCPProxy) error {
	host.LastPacketDate = time.Now()
	return nil
}

func (host *Host) StartProxies() {
	// for _, tcpProxy := range host.TCPProxies {
	// 	// tcpProxy.Start()
	// }
}

func (host *Host) Dispose() {
	for _, tcpProxy := range host.TCPProxies {
		tcpProxy.Stop()

		log.Printf("TCP Proxy %s disposed\n", tcpProxy.ListenAddr)
	}

	log.Println("Waiting for all proxies to stop")

	host.waitGroup.Wait()

	log.Println("Host disposed : " + host.Config.Name)
}
