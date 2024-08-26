package host

import (
	"context"
	"log"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/proxies"
	"time"
)

type Host struct {
	TCPProxies     map[string]*proxies.TCPProxy
	UDPProxies     map[string]*proxies.UDPProxy
	Started        bool
	LastPacketDate time.Time
	Config         *config.HostConfig
}

func NewHost(hostConfig *config.HostConfig) *Host {
	h := &Host{
		Config: hostConfig,
	}

	for _, proxyConfig := range hostConfig.Proxies {
		if proxyConfig.Protocol == "tcp" {
			h.TCPProxies[proxyConfig.Name] = proxies.NewTCPProxy(context.TODO(), &proxies.TCPProxyArgs{
				HostConfig:     hostConfig,
				ProxyConfig:    proxyConfig,
				HostStarted:    h.HostStarted,
				StartHost:      h.StartHost,
				PacketReceived: h.PacketReceived,
			})
		} else if proxyConfig.Protocol == "udp" {

		}
	}

	// h.StartHost()
	h.Started = true
	h.LastPacketDate = time.Now()

	return h
}

func (h *Host) HostStarted(proxy *proxies.TCPProxy) (bool, error) {
	return h.Started, nil
}

func (h *Host) StartHost(proxy *proxies.TCPProxy) error {
	log.Println("Starting host : " + h.Config.Name)
	return nil
}

func (h *Host) StopHost() error {
	log.Println("Stopping host : " + h.Config.Name)
	return nil
}

func (h *Host) PacketReceived(proxy *proxies.TCPProxy) error {
	h.LastPacketDate = time.Now()
	return nil
}

func (h *Host) StartProxies() {
	// for _, tcpProxy := range h.TCPProxies {
	// 	// tcpProxy.Start()
	// }
}
