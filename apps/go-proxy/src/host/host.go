package host

import (
	"context"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/docker"
	"mgarnier11/go-proxy/proxies"
	"slices"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type Host struct {
	Proxies        map[string]*proxies.TCPProxy
	Started        bool
	LastPacketDate time.Time
	Config         *config.HostConfig

	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewHost(hostConfig *config.HostConfig) *Host {
	ctx, cancel := context.WithCancel(context.Background())

	host := &Host{
		Config:    hostConfig,
		Proxies:   make(map[string]*proxies.TCPProxy),
		waitGroup: sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
	}

	host.StartHost(nil)
	host.Started = true
	host.LastPacketDate = time.Now()

	host.setupProxies(hostConfig.Proxies)

	host.waitGroup.Add(1)
	go func() {
		defer host.waitGroup.Done()
		for dockerProxies := range docker.SetupDockerContainersListener(host.ctx, hostConfig.SSHUsername, hostConfig.Ip) {
			log.Infof("%-10s received %d docker proxies", host.Config.Name, len(dockerProxies))
			host.setupProxies(dockerProxies)
		}

		log.Infof("%-10s stopped listening for docker containers", host.Config.Name)
	}()

	log.Infof("%-10s created", host.Config.Name)

	return host
}

func (host *Host) setupProxies(proxyConfigs []*config.ProxyConfig) {
	for name := range host.Proxies {
		exists := slices.ContainsFunc(proxyConfigs, func(proxy *config.ProxyConfig) bool {
			return proxy.Name == name
		})

		if !exists {
			host.DisposeProxy(name)
		}
	}

	for _, proxyConfig := range proxyConfigs {
		if host.Proxies[proxyConfig.Name] != nil {
			log.Debugf("%-10s %-20s already exists", host.Config.Name, proxyConfig.Name)
			continue
		}

		host.Proxies[proxyConfig.Name] = proxies.NewTCPProxy(&proxies.TCPProxyArgs{
			HostName:       host.Config.Name,
			HostIp:         host.Config.Ip,
			ProxyConfig:    proxyConfig,
			HostStarted:    host.HostStarted,
			StartHost:      host.StartHost,
			PacketReceived: host.PacketReceived,
		})

		go host.Proxies[proxyConfig.Name].Start(&host.waitGroup)
	}
}

func (host *Host) HostStarted() (bool, error) {
	return host.Started, nil
}

func (host *Host) StartHost(proxy *proxies.TCPProxy) error {
	log.Infof("%-10s Starting", host.Config.Name)
	return nil
}

func (host *Host) StopHost() error {
	log.Infof("%-10s Stopping", host.Config.Name)
	return nil
}

func (host *Host) PacketReceived(proxy *proxies.TCPProxy) error {
	host.LastPacketDate = time.Now()
	return nil
}

func (host *Host) DisposeProxy(proxyName string) {
	proxy := host.Proxies[proxyName]

	if proxy == nil {
		log.Errorf("%-10s %s: proxy does not exist", host.Config.Name, proxyName)
		return
	}

	proxy.Stop()
	delete(host.Proxies, proxyName)

	log.Infof("%-10s %s: disposed", host.Config.Name, proxyName)
}

func (host *Host) Dispose() {
	log.Infof("%-10s disposing", host.Config.Name)

	host.cancel()

	for name := range host.Proxies {
		host.DisposeProxy(name)
	}

	host.waitGroup.Wait()

	log.Infof("%-10s disposed", host.Config.Name)
}
