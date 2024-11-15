package host

import (
	"context"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/docker"
	"mgarnier11/go-proxy/proxies"
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
			log.Infof("Host %s received %d docker proxies", host.Config.Name, len(dockerProxies))
			host.setupProxies(dockerProxies)
		}

		log.Infof("Host %s stopped listening for docker containers", host.Config.Name)
	}()

	log.Infof("Host %s created", host.Config.Name)

	return host
}

func (host *Host) setupProxies(proxyConfigs []*config.ProxyConfig) {
	for _, proxyConfig := range proxyConfigs {
		if host.Proxies[proxyConfig.Name] != nil {
			log.Debugf("Proxy %s already exists", proxyConfig.Name)
			continue
		}

		host.Proxies[proxyConfig.Name] = proxies.NewTCPProxy(&proxies.TCPProxyArgs{
			HostIp:         host.Config.Ip,
			ProxyConfig:    proxyConfig,
			HostStarted:    host.HostStarted,
			StartHost:      host.StartHost,
			PacketReceived: host.PacketReceived,
		})

		host.waitGroup.Add(1)
		go host.Proxies[proxyConfig.Name].Start(&host.waitGroup)
	}
}

func (host *Host) Ref(proxy *proxies.TCPProxy) (bool, error) {
	return host.Started, nil
}

func (host *Host) HostStarted(proxy *proxies.TCPProxy) (bool, error) {
	return host.Started, nil
}

func (host *Host) StartHost(proxy *proxies.TCPProxy) error {
	log.Infof("Starting host : %s", host.Config.Name)
	return nil
}

func (host *Host) StopHost() error {
	log.Infof("Stopping host : %s", host.Config.Name)
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

func (host *Host) DisposeProxy(proxyName string) {
	proxy := host.Proxies[proxyName]

	if proxy == nil {
		log.Errorf("Cant dispose, proxy %s does not exist", proxyName)
		return
	}

	proxy.Stop()
	delete(host.Proxies, proxyName)

	log.Infof("Proxy %s disposed", proxy.ListenAddr)
}

func (host *Host) Dispose() {
	log.Infof("Disposing host %s", host.Config.Name)

	host.cancel()

	for name := range host.Proxies {
		host.DisposeProxy(name)
	}

	log.Infof("Host %s waiting for proxies to stop", host.Config.Name)

	host.waitGroup.Wait()

	log.Infof("Host %s disposed", host.Config.Name)
}
