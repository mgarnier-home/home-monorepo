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

	stopping bool
	starting bool

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

	dockerProxies, _ := docker.GetProxiesFromDocker(hostConfig.SSHUsername, hostConfig.Ip)

	host.setupProxies(slices.Concat(dockerProxies, hostConfig.Proxies))

	go host.setupContainersListener()

	log.Infof("%-10s created", host.Config.Name)

	return host
}

func (host *Host) setupContainersListener() {
	host.waitGroup.Add(1)
	log.Infof("%-10s listening for docker containers", host.Config.Name)

	defer func() {
		log.Infof("%-10s stopped listening for docker containers", host.Config.Name)
		host.waitGroup.Done()
	}()

	ticker := time.NewTicker(50 * time.Second)

	for {
		select {
		case <-ticker.C:
			log.Infof("%-10s started: %v", host.Config.Name, host.Started)
			if host.Started {
				log.Infof("Getting proxies from docker")

				proxies, err := docker.GetProxiesFromDocker(host.Config.SSHUsername, host.Config.Ip)

				if err != nil {
					log.Errorf("%-10s failed to get proxies from docker: %v", host.Config.Name, err)
				} else {
					host.setupProxies(proxies)
				}
			}
		case <-host.ctx.Done():
			// Ensure we break out of the loop if the context is cancelled
			ticker.Stop()
			return
		}
	}

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

func (host *Host) StopHost() {
	if host.stopping {
		log.Infof("%-10s Already stopping", host.Config.Name)

		return
	}

	host.stopping = true
	defer func() {
		host.stopping = false
	}()

	log.Infof("%-10s Stopping", host.Config.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	go func() {

		err := sendSSHCommand(ctx, host.Config, "sudo pm-suspend &")

		if err != nil {
			log.Errorf("%-10s failed to stop host: %v", host.Config.Name, err)
		}
	}()

	for {
		select {
		case <-ctx.Done(): // Context timeout or cancellation
			log.Warnf("%-10s Context canceled or timed out", host.Config.Name)
			return
		default:
			hostPinged, err := getHostStatus(host.Config.Ip)

			log.Infof("%-10s Host pinged: %v", host.Config.Name, hostPinged)

			if err != nil {
				log.Errorf("%-10s failed to check host status: %v", host.Config.Name, err)
			} else if !hostPinged {
				log.Infof("%-10s Stopped", host.Config.Name)

				host.Started = false

				return
			}

			time.Sleep(1 * time.Second)
		}
	}

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
