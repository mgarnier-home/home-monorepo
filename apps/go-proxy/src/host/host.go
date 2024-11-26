package host

import (
	"context"
	"fmt"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/docker"
	"mgarnier11/go-proxy/hostState"
	"mgarnier11/go-proxy/proxies"
	"mgarnier11/go/logger"
	"slices"
	"strings"
	"sync"
	"time"
)

type Host struct {
	Proxies        map[string]*proxies.TCPProxy
	State          hostState.State
	LastPacketDate time.Time
	Config         *config.HostConfig

	logger *logger.Logger

	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewHost(hostConfig *config.HostConfig) *Host {
	ctx, cancel := context.WithCancel(context.Background())

	host := &Host{
		Config:    hostConfig,
		Proxies:   make(map[string]*proxies.TCPProxy),
		State:     hostState.Stopped,
		waitGroup: sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
		logger:    logger.NewLogger(fmt.Sprintf("[%s]", strings.ToUpper(hostConfig.Name)), "%-10s ", nil),
	}

	host.logger.Infof("created")

	go host.setupHostLoop()

	host.StartHost()
	host.LastPacketDate = time.Now()

	dockerProxies, _ := docker.GetProxiesFromDocker(hostConfig.SSHUsername, hostConfig.Ip, host.logger)

	host.setupProxies(slices.Concat(dockerProxies, hostConfig.Proxies))

	return host
}

func (host *Host) setupHostLoop() {
	host.waitGroup.Add(1)
	host.logger.Infof("starting host loop")

	defer func() {
		host.logger.Infof("stopping host loop")
		host.waitGroup.Done()
	}()

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			host.updateState()

			// if host.State == Started {
			// 	host.logger.Infof("Getting proxies from docker")

			// 	proxies, err := docker.GetProxiesFromDocker(host.Config.SSHUsername, host.Config.Ip, host.logger)

			// 	if err != nil {
			// 		host.logger.Errorf("failed to get proxies from docker: %v", err)
			// 	} else {
			// 		host.setupProxies(proxies)
			// 	}
			// }
		case <-host.ctx.Done():
			// Ensure we break out of the loop if the context is cancelled
			ticker.Stop()
			return
		}
	}
}

func (host *Host) updateState() {
	pingSuccess, err := getHostStatus(host.Config.Ip)

	if err != nil {
		host.logger.Errorf("failed to check host status: %v", err)
		return
	}

	if host.State == hostState.Started && !pingSuccess {
		host.State = hostState.Stopped
	} else if host.State == hostState.Stopped && pingSuccess {
		host.State = hostState.Started
	} else if host.State == hostState.Starting && pingSuccess {
		host.State = hostState.Started
	} else if host.State == hostState.Stopping && !pingSuccess {
		host.State = hostState.Stopped
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
			host.logger.Debugf("%s already exists", proxyConfig.Name)
			continue
		}

		host.Proxies[proxyConfig.Name] = proxies.NewTCPProxy(&proxies.TCPProxyArgs{
			HostIp:         host.Config.Ip,
			ProxyConfig:    proxyConfig,
			HostState:      &host.State,
			StartHost:      host.StartHost,
			PacketReceived: host.PacketReceived,
		}, host.logger)

		go host.Proxies[proxyConfig.Name].Start(&host.waitGroup)
	}
}

func (host *Host) StartHost() error {
	if host.State != hostState.Stopped {
		host.logger.Infof("Cannot start host, state is not stopped : %s", host.State.String())
		return nil
	}

	host.State = hostState.Starting

	if packet, err := newMagicPacket(host.Config.MacAddress); err == nil {
		packet.send("255.255.255.255")
		host.logger.Debugf("Sent magic packet to start host")
	} else {
		return fmt.Errorf("failed to send magic packet: %v", err)
	}

	i := 0
	for host.State == hostState.Starting && i < 20 {
		time.Sleep(1 * time.Second)
		i++
	}

	if host.State == hostState.Starting || i >= 20 {
		host.State = hostState.Stopped
		return fmt.Errorf("Host took too long to start")
	} else {
		return nil
	}
}

func (host *Host) StopHost() {
	if host.State != hostState.Started {
		host.logger.Infof("Cannot stop host, state is not started : %s", host.State.String())
		return
	}

	host.State = hostState.Stopping

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := sendSSHCommand(ctx, host.Config, "sudo pm-suspend &")

		if err != nil {
			host.logger.Errorf("failed to stop host: %v", err)
		}
	}()

	i := 0

	for host.State == hostState.Stopping && i < 20 {
		time.Sleep(1 * time.Second)
		i++
	}

	if host.State == hostState.Stopping || i >= 20 {
		host.State = hostState.Started
		host.logger.Errorf("Host took too long to stop")
	}
}

func (host *Host) PacketReceived(proxy *proxies.TCPProxy) error {
	host.LastPacketDate = time.Now()
	return nil
}

func (host *Host) DisposeProxy(proxyName string) {
	proxy := host.Proxies[proxyName]

	if proxy == nil {
		host.logger.Errorf("%s: proxy does not exist", proxyName)
		return
	}

	proxy.Stop()
	delete(host.Proxies, proxyName)

	host.logger.Infof("%s: disposed", proxyName)
}

func (host *Host) Dispose() {
	host.logger.Infof("disposing")

	host.cancel()

	for name := range host.Proxies {
		host.DisposeProxy(name)
	}

	host.waitGroup.Wait()

	host.logger.Infof("disposed")
}
